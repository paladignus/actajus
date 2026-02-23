// Package identity
package identity

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/infrastructure/security"
	"github.com/paladignus/actajus/internal/module/identity/presentation/grpc/handler"
	"github.com/paladignus/actajus/internal/shared/infrastructure/clock"
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

// type JWTConfig struct {
// 	Issuer       string
// 	Audience     string
// 	AccessTTL    time.Duration
// 	AccessSecret string // usado apenas pela infra do HS256
// }
//
// type SessionConfig struct {
// 	RefreshTTL  time.Duration
// 	MaxSessions int
// }
//
// type PasswordResetConfig struct {
// 	ResetTTL time.Duration
// }

type Module struct {
	Handler handler.AuthHandler
	// Exponho ValidateAccess pra você plugar no AuthInterceptor no main.
	ValidateAccess usecase.ValidateAccess
}

func (m Module) Route(opts ...connect.HandlerOption) (string, http.Handler) {
	return identityv1connect.NewAuthServiceHandler(m.Handler, opts...)
}

func NewModule() (Module, error) {
	// infra services puros (sem DB)
	clk := clock.NewSystemClock()
	refreshSvc := security.NewRefreshTokenService()
	hasher := security.NewArgon2idPasswordHasher()
	accessSvc, err := security.NewHS256AccessTokenService(
		dep.JWTConfig.AccessSecret,
		dep.JWTConfig.Issuer,
		dep.JWTConfig.Audience,
	)
	if err != nil {
		return Module{}, err
	}
	return Module{}, err
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
