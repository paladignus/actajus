// Package identity
package identity

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/redis/go-redis/v9"

	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/cache"
	identitypg "github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
	identityhandler "github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
	identityrbac "github.com/paladignus/actajus/internal/module/identity/presentation/rbac"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
	sharedsvc "github.com/paladignus/actajus/internal/shared/application/service"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/clock"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	postgresShared "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/presentation/interceptor"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
	identityv1connect "github.com/paladignus/actajus/proto/identity/v1/identityv1connect"
)

type Module struct {
	ValidateAccess usecase.ValidateAccess
	RBACChecker    interceptor.PermissionChecker
	authImpl       *identityhandler.AuthHandler
	rbacAdminImpl  *identityhandler.RbacAdminHandler
	logger         sharedrepo.Logger
	db             postgresShared.Executor
	rdb            redis.UniversalClient
}

type Dependencies struct {
	Logger          sharedrepo.Logger
	DB              postgresShared.Executor
	RDB             redis.UniversalClient
	Config          config.AuthConfig
	UoW             uow.UnitOfWork
	Repository      repository.Factory
	OutboxFactory   messaging.OutboxFactory
	Serializer      sharedsvc.MessageSerializer
	IDGenerator     sharedsvc.IDGenerator
	SkipRBACRebuild bool
}

func (m Module) Mount(mux *http.ServeMux, opts ...connect.HandlerOption) {
	mux.Handle(identityv1connect.NewAuthServiceHandler(m.authImpl, opts...))
	mux.Handle(identityv1connect.NewRbacAdminServiceHandler(m.rbacAdminImpl, opts...))
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
	// Repository com cache fallback para consultas de usuários por role
	roleUserQueryRepo := cache.NewCachedRoleUserQueryRepository(
		roleUserAdminRepo,
		roleUsersIndex,
		dep.Logger,
	)
	loginUC := usecase.NewLogin(
		dep.Repository.User(),
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
		dep.Repository.User(),
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
		dep.Repository.User(),
		cachedSessionRepo,
		hasher,
		clk,
		authMapper,
		true,
	)
	requestResetUC := usecase.NewRequestPasswordReset(
		dep.UoW,
		dep.Repository,
		dep.OutboxFactory,
		refreshSvc,
		clk,
		dep.Config.PasswordResetConfig,
		*authMapper,
		dep.Serializer,
		dep.IDGenerator,
		true,
	)
	confirmResetUC := usecase.NewConfirmPasswordReset(
		dep.Repository.User(),
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
		roleUserQueryRepo,
		permRoleAdminRepo,
		authzSvc,
		rbacMapper,
	)
	revokeUC := usecase.NewRevokePermissionFromRole(
		roleUserQueryRepo,
		permRoleAdminRepo,
		authzSvc,
		rbacMapper,
	)
	rbacAdminImpl := identityhandler.NewRbacAdminHandler(
		assignUC,
		removeUC,
		grantUC,
		revokeUC,
	)
	rbacChecker := identityrbac.NewChecker(authzSvc)
	module := Module{
		ValidateAccess: validateAccessUC,
		RBACChecker:    rbacChecker,
		authImpl:       &authImpl,
		rbacAdminImpl:  rbacAdminImpl,
		logger:         dep.Logger,
		db:             dep.DB,
		rdb:            dep.RDB,
	}

	// Rebuild RBAC indexes automaticamente (a menos que seja explicitamente desabilitado)
	if !dep.SkipRBACRebuild {
		if err := module.rebuildRBACIndexes(context.Background()); err != nil {
			dep.Logger.Error(context.Background(), "⚠️ failed to rebuild RBAC indexes", "error", err)
		} else {
			dep.Logger.Info(context.Background(), "✅ RBAC indexes rebuilt successfully")
		}
	}
	return module, nil
}

// rebuildRBACIndexes reconstrói os índices RBAC no Redis a partir dos dados do Postgres.
// Este método é chamado automaticamente durante a inicialização do módulo.
func (m Module) rebuildRBACIndexes(ctx context.Context) error {
	pairs, err := identitypg.ListAllRoleUsers(ctx, m.db)
	if err != nil {
		return err
	}
	roleUsersIndex := security.NewRBACRoleUsersIndex(
		m.rdb,
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

// RebuildRBACIndexes reconstrói os índices RBAC no Redis.
// Deprecated: Use o rebuild automático via NewModule ou chame diretamente no módulo.
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
