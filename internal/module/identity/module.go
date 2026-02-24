// Package identity
package identity

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
	"github.com/paladignus/actajus/proto/identity/v1/identityv1connect"
)

// Dependencies são as dependências que vêm de fora (infra repos etc).
// Assim você consegue montar o módulo em k8s com Postgres/Redis depois,
// sem mudar nada na aplicação.
// type Dependencies struct {
// 	logger        repository.Logger
// 	users         repository.UserRepository
// 	Sessions      repository.SessionRepository
// 	PasswordReset repository.PasswordResetRepository
// 	JWTConfig     config.JWTConfig
// 	SessionConfig config.SessionConfig
// 	ResetConfig   config.PasswordResetConfig
// }

type Module struct {
	Handler handler.AuthHandler
}

func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
	return identityv1connect.NewAuthServiceHandler(m.Handler, opts...)
}

func NewModule(
	pool *pgxpool.Pool,
	logger repository.Logger,
) Module {
	// h := handler.NewPersonHandler(usecase, logger)
	return Module{}
}

// mapper (validação/normalização)
// Aqui eu assumo que você já tem um validator injetável no mapper.
// Se o seu AuthMapper ainda cria validator internamente, ok também.
// authMapper := mapper.NewAuthMapper() // ajuste pro seu construtor real
// usecases
// loginUC := usecase.NewLogin(
// 	dep.Users,
// 	dep.Sessions,
// 	hasher,
// 	accessSvc,
// 	refreshSvc,
// 	service.JWTConfig{
// 		Issuer:    dep.JWTConfig.Issuer,
// 		Audience:  dep.JWTConfig.Audience,
// 		AccessTTL: dep.JWTConfig.AccessTTL,
// 	},
// 	service.SessionConfig{
// 		RefreshTTL:  dep.SessionConfig.RefreshTTL,
// 		MaxSessions: dep.SessionConfig.MaxSessions,
// 	},
// 	clk,
// 	authMapper,
// )
//
// refreshUC := usecase.NewRefresh(
// 	dep.Sessions,
// 	accessSvc,
// 	refreshSvc,
// 	service.JWTConfig{
// 		Issuer:    dep.JWTConfig.Issuer,
// 		Audience:  dep.JWTConfig.Audience,
// 		AccessTTL: dep.JWTConfig.AccessTTL,
// 	},
// 	service.SessionConfig{
// 		RefreshTTL:  dep.SessionConfig.RefreshTTL,
// 		MaxSessions: dep.SessionConfig.MaxSessions,
// 	},
// 	clk,
// 	authMapper,
// 	true, // revoke session on refresh mismatch (você pediu revogar automaticamente)
// )
//
// logoutUC := usecase.NewLogout(dep.Sessions, clk, authMapper)
// logoutAllUC := usecase.NewLogoutAll(dep.Sessions, clk, authMapper)
// changePasswordUC := usecase.NewChangePassword(
// 	dep.Users,
// 	dep.Sessions,
// 	hasher,
// 	clk,
// 	authMapper,
// 	true, // revoke all sessions on password change
// )
//
// requestResetUC := usecase.NewRequestPasswordReset(
// 	dep.Users,
// 	dep.PasswordReset,
// 	refreshSvc,
// 	clk,
// 	usecase.PasswordResetConfig{ResetTTL: dep.ResetConfig.ResetTTL},
// 	authMapper,
// 	true, // revoke previous reset tokens
// )
//
// confirmResetUC := usecase.NewConfirmPasswordReset(
// 	dep.Users,
// 	dep.Sessions,
// 	dep.PasswordReset,
// 	hasher,
// 	refreshSvc,
// 	clk,
// 	authMapper,
// 	true, // revoke all sessions on reset confirm
// )
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
