// Package identity
package identity

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/redis/go-redis/v9"

	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/cache"
	identitypg "github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
	identityhandler "github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
	identityrbac "github.com/paladignus/actajus/internal/module/identity/presentation/rbac"
	sharedrepo "github.com/paladignus/actajus/internal/shared/domain/repository"
	"github.com/paladignus/actajus/internal/shared/infrastructure/clock"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
	identityv1connect "github.com/paladignus/actajus/proto/identity/v1/identityv1connect"
)

type Dependencies struct {
	Logger sharedrepo.Logger
	DB     postgresShared.PgxPool
	RDB    redis.UniversalClient
	Users  identityrepo.UserRepository
	Config config.AuthConfig
}

type Module struct {
	ValidateAccess usecase.ValidateAccess
	RBACChecker    interceptor.PermissionChecker
	authImpl       *identityhandler.AuthHandler
	rbacAdminImpl  *identityhandler.RbacAdminHandler
}

func (m Module) Mount(mux *http.ServeMux, opts ...connect.HandlerOption) {
	authPath, authHTTPHandler := identityv1connect.NewAuthServiceHandler(m.authImpl, opts...)
	mux.Handle(authPath, authHTTPHandler)
	rbacAdminPath, rbacAdminHTTPHandler := identityv1connect.NewRbacAdminServiceHandler(m.rbacAdminImpl, opts...)
	mux.Handle(rbacAdminPath, rbacAdminHTTPHandler)
}

func NewModule(dep Dependencies) (Module, error) {
	clk := clock.NewSystemClock()
	authValidator := validation.New()
	projection := mapper.NewAuthProjectionMapper()
	authMapper := mapper.NewAuthMapper(authValidator)
	rbacValidator := validation.New()
	rbacMapper := mapper.NewRBACAdminMapper(rbacValidator)
	refreshSvc := security.NewRefreshTokenService()
	hasher := security.NewArgon2idPasswordHasher()
	accessSvc, err := security.NewHS256AccessTokenService(
		dep.Config.AccessSecret,
		dep.Config.Issuer,
		dep.Config.Audience,
	)
	if err != nil {
		return Module{}, err
	}
	pgSessionRepo := identitypg.NewSession(dep.DB)
	cachedSessionRepo := cache.NewCachedSession(
		pgSessionRepo,
		dep.RDB,
		dep.Logger,
		cache.WithPrefix("actajus:"),
		cache.WithFallbackToPostgres(true),
	)
	passwordResetRepo := identitypg.NewPasswordReset(dep.DB)
	authzRepo := identitypg.NewAuthorization(dep.DB)
	authzSvc := security.NewAuthorizationService(
		authzRepo,
		dep.RDB,
		security.WithAuthzPrefix("rbac:"),
		security.WithAuthzTTL(10*time.Minute),
	)
	roleUsersIndex := security.NewRBACRoleUsersIndex(
		dep.RDB,
		security.WithRoleUsersIndexPrefix("rbac:"),
	)
	roleUserAdminRepo := identitypg.NewRoleUserAdminRepository(dep.DB)
	permRoleAdminRepo := identitypg.NewPermissionRoleAdminRepository(dep.DB)
	loginUC := usecase.NewLogin(
		dep.Users,
		cachedSessionRepo,
		hasher,
		refreshSvc,
		accessSvc,
		clk,
		dep.Config,
		*authMapper,
		*projection,
	)

	refreshUC := usecase.NewRefresh(
		cachedSessionRepo,
		dep.Users,
		refreshSvc,
		accessSvc,
		clk,
		dep.Config,
		*authMapper,
		*projection,
	)

	logoutUC := usecase.NewLogout(cachedSessionRepo, *authMapper)
	logoutAllUC := usecase.NewLogoutAll(cachedSessionRepo, *authMapper)

	changePasswordUC := usecase.NewChangePassword(
		dep.Users,
		cachedSessionRepo,
		hasher,
		clk,
		authMapper,
		true,
	)

	requestResetUC := usecase.NewRequestPasswordReset(
		dep.Users,
		passwordResetRepo,
		refreshSvc,
		clk,
		dep.Config.PasswordResetConfig,
		// usecase.PasswordResetConfig{
		// 	ResetTTL: dep.PasswordResetTTL.ResetTTL,
		// },
		*authMapper,
		true,
	)

	confirmResetUC := usecase.NewConfirmPasswordReset(
		dep.Users,
		cachedSessionRepo,
		passwordResetRepo,
		hasher,
		refreshSvc,
		clk,
		*authMapper,
		true,
	)

	validateAccessUC := usecase.NewValidateAccess(
		accessSvc,
		cachedSessionRepo,
		clk,
		true,
	)

	authImpl := identityhandler.NewAuthHandler(
		loginUC,
		refreshUC,
		logoutUC,
		logoutAllUC,
		changePasswordUC,
		requestResetUC,
		confirmResetUC,
	)

	// rbac admin usecases
	assignUC := usecase.NewAssignRoleToUser(
		roleUserAdminRepo,
		authzSvc,
		roleUsersIndex,
		rbacMapper,
	)

	removeUC := usecase.NewRemoveRoleFromUser(
		roleUserAdminRepo,
		authzSvc,
		roleUsersIndex,
		rbacMapper,
	)

	grantUC := usecase.NewGrantPermissionToRole(
		roleUserAdminRepo,
		permRoleAdminRepo,
		authzSvc,
		roleUsersIndex,
		rbacMapper,
	)

	revokeUC := usecase.NewRevokePermissionFromRole(
		roleUserAdminRepo,
		permRoleAdminRepo,
		authzSvc,
		roleUsersIndex,
		rbacMapper,
	)

	rbacAdminImpl := identityhandler.NewRbacAdminHandler(
		assignUC,
		removeUC,
		grantUC,
		revokeUC,
	)

	rbacChecker := identityrbac.NewChecker(authzSvc)

	return Module{
		ValidateAccess: validateAccessUC,
		RBACChecker:    rbacChecker,
		authImpl:       &authImpl,
		rbacAdminImpl:  rbacAdminImpl,
	}, nil
}

func RebuildRBACIndexes(ctx context.Context, dep Dependencies) error {
	pairs, err := identitypg.ListAllRoleUsers(ctx, dep.DB)
	if err != nil {
		return err
	}

	roleUsersIndex := security.NewRBACRoleUsersIndex(
		dep.RDB,
		security.WithRoleUsersIndexPrefix("rbac:"),
	)

	rolePairs := make([]security.RoleUserPair, 0, len(pairs))
	for _, p := range pairs {
		rolePairs = append(rolePairs, security.RoleUserPair{
			IDRole: p.IDRole,
			IDUser: p.IDUser,
		})
	}

	return security.RebuildRoleUsersIndex(ctx, roleUsersIndex, rolePairs)
}

// func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
// 	return identityv1connect.NewAuthServiceHandler(m.Handler, opts...)
// }

// func NewModule(
// 	pool *pgxpool.Pool,
// 	logger repository.Logger,
// 	config config.AuthConfig,
// 	rdb *redis.Client,
// ) (Module, error) {
// 	clk := clock.NewSystemClock()
// 	refresh := security.NewRefreshTokenService()
// 	hasher := security.NewArgon2idPasswordHasher()
// 	accessSvc, err := security.NewHS256AccessTokenService(
// 		config.AccessSecret,
// 		config.Issuer,
// 		config.Audience,
// 	)
// 	if err != nil {
// 		return Module{}, err
// 	}
// 	// validation := validation.New(nil))
// 	v := validation.New()
// 	projection := mapper.NewAuthProjectionMapper()
// 	mapper := mapper.NewAuthMapper(v)
// 	userRepo := postgres.NewUser(pool)
// 	sessionRepo := postgres.NewSession(pool)
// 	cacheRepo := cache.NewCachedSession(
// 		sessionRepo,
// 		rdb,
// 		logger,
// 		cache.WithPrefix("actajus:"),
// 		cache.WithFallbackToPostgress(true),
// 	)
// 	passResetRepo := postgres.NewPasswordReset(pool)
// 	loginUC := usecase.NewLogin(
// 		userRepo,
// 		cacheRepo,
// 		hasher,
// 		refresh,
// 		accessSvc,
// 		clk,
// 		config,
// 		*mapper,
// 		*projection,
// 	)
// 	refreshUC := usecase.NewRefresh(
// 		cacheRepo,
// 		userRepo,
// 		refresh,
// 		accessSvc,
// 		clk,
// 		config,
// 		*mapper,
// 		*projection,
// 	)
// 	logoutUC := usecase.NewLogout(&sessionRepo, *mapper)
// 	logoutAllUC := usecase.NewLogoutAll(&sessionRepo, *mapper)
// 	changePasswordUC := usecase.NewChangePassword(
// 		userRepo,
// 		&sessionRepo,
// 		hasher,
// 		clk,
// 		mapper,
// 		true,
// 	)
// 	requestResetUC := usecase.NewRequestPasswordReset(
// 		userRepo,
// 		passResetRepo,
// 		refresh,
// 		clk,
// 		config.PasswordResetConfig,
// 		*mapper,
// 		true,
// 	)
// 	confirmResetUC := usecase.NewConfirmPasswordReset(
// 		userRepo,
// 		cacheRepo,
// 		passResetRepo,
// 		hasher,
// 		refresh,
// 		clk,
// 		*mapper,
// 		true, // revoke all sessions on reset confirm
// 	)
// 	validateAccessUC := usecase.NewValidateAccess(
// 		accessSvc,
// 		cacheRepo,
// 		clk,
// 		// *mapper,
// 		true, // checkSession=true (bom para revogação)
// 	)
// 	h := handler.NewAuthHandler(
// 		loginUC,
// 		refreshUC,
// 		logoutUC,
// 		logoutAllUC,
// 		changePasswordUC,
// 		requestResetUC,
// 		confirmResetUC,
// 	)
// 	return Module{h, validateAccessUC}, nil
// }
