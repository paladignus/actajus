// Package identity
package identity

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/cache"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
	"github.com/paladignus/actajus/internal/shared/infrastructure/clock"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
	"github.com/paladignus/actajus/proto/identity/v1/identityv1connect"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	Handler        handler.AuthHandler
	ValidateAccess usecase.ValidateAccess
}

func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
	return identityv1connect.NewAuthServiceHandler(m.Handler, opts...)
}

func NewModule(
	pool *pgxpool.Pool,
	logger repository.Logger,
	config config.AuthConfig,
	rdb *redis.Client,
) (Module, error) {
	clk := clock.NewSystemClock()
	refresh := security.NewRefreshTokenService()
	hasher := security.NewArgon2idPasswordHasher()
	accessSvc, err := security.NewHS256AccessTokenService(
		config.AccessSecret,
		config.Issuer,
		config.Audience,
	)
	if err != nil {
		return Module{}, err
	}
	// validation := validation.New(nil))
	v := validation.New()
	projection := mapper.NewAuthProjectionMapper()
	mapper := mapper.NewAuthMapper(v)
	userRepo := postgres.NewUser(pool)
	sessionRepo := postgres.NewSession(pool)
	cacheRepo := cache.NewCachedSession(
		sessionRepo,
		rdb,
		cache.WithPrefix("actajus:"),
		cache.WithFallbackToPostgress(true),
	)
	passResetRepo := postgres.NewPasswordReset(pool)
	loginUC := usecase.NewLogin(
		userRepo,
		cacheRepo,
		hasher,
		refresh,
		accessSvc,
		clk,
		config,
		*mapper,
		*projection,
	)
	refreshUC := usecase.NewRefresh(
		cacheRepo,
		userRepo,
		refresh,
		accessSvc,
		clk,
		config,
		*mapper,
		*projection,
	)
	logoutUC := usecase.NewLogout(sessionRepo, *mapper)
	logoutAllUC := usecase.NewLogoutAll(sessionRepo, *mapper)
	changePasswordUC := usecase.NewChangePassword(
		userRepo,
		sessionRepo,
		hasher,
		clk,
		mapper,
		true,
	)
	requestResetUC := usecase.NewRequestPasswordReset(
		userRepo,
		passResetRepo,
		refresh,
		clk,
		config.PasswordResetConfig,
		*mapper,
		true,
	)
	confirmResetUC := usecase.NewConfirmPasswordReset(
		userRepo,
		cacheRepo,
		passResetRepo,
		hasher,
		refresh,
		clk,
		*mapper,
		true, // revoke all sessions on reset confirm
	)
	validateAccessUC := usecase.NewValidateAccess(
		accessSvc,
		cacheRepo,
		clk,
		*mapper,
		true, // checkSession=true (bom para revogação)
	)
	h := handler.NewAuthHandler(
		loginUC,
		refreshUC,
		logoutUC,
		logoutAllUC,
		changePasswordUC,
		requestResetUC,
		confirmResetUC,
	)
	return Module{h, validateAccessUC}, nil
}

//
// validateAccessUC := usecase.NewValidateAccess(
// 	accessSvc,
// 	dep.Sessions,
// 	clk,
// 	authMapper,
// 	true, // checkSession=true (bom para revogação)
// )
//
// // handler Connect-Go (impl do service)
// h := handler.NewAuthHandler(
// 	loginUC,
// 	refreshUC,
// 	logoutUC,
// 	logoutAllUC,
// 	changePasswordUC,
// 	requestResetUC,
// 	confirmResetUC,
// )

// return Module{
// 	Handler:        *h,
// 	ValidateAccess: validateAccessUC,
// }, nil
// }
