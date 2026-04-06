package auth

import (
	"errors"
	"net/http"

	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
)

func LoginError(err error) (int, string) {
	switch {
	case errors.Is(err, identitydomain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Email ou senha invalidos."
	case errors.Is(err, identitydomain.ErrEmailNotVerified):
		return http.StatusForbidden, "Seu email ainda nao foi verificado."
	case errors.Is(err, identitydomain.ErrUserBlocked):
		return http.StatusForbidden, "Usuario bloqueado."
	case errors.Is(err, identitydomain.ErrSessionLimit):
		return http.StatusTooManyRequests, "Limite de sessoes ativas atingido."
	default:
		return http.StatusBadRequest, "Nao foi possivel autenticar."
	}
}

func RegisterError(err error) (int, string) {
	switch {
	default:
		return http.StatusBadRequest, "Se o cadastro puder ser concluido, voce recebera um email com os proximos passos."
	}
}

func VerifyEmailError(err error) (int, string) {
	switch {
	case errors.Is(err, identitydomain.ErrEmailAlreadyVerified):
		return http.StatusOK, "Este email ja foi verificado anteriormente."
	case errors.Is(err, identitydomain.ErrEmailVerificationExpired):
		return http.StatusGone, "O link de verificacao expirou."
	case errors.Is(err, identitydomain.ErrEmailVerificationUsed):
		return http.StatusGone, "Este link de verificacao ja foi utilizado."
	case errors.Is(err, identitydomain.ErrEmailVerificationInvalid):
		return http.StatusBadRequest, "O token de verificacao e invalido."
	default:
		return http.StatusBadRequest, "Nao foi possivel validar seu email."
	}
}
