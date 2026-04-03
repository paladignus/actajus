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
	case errors.Is(err, identitydomain.ErrUserBlocked):
		return http.StatusForbidden, "Usuario bloqueado."
	case errors.Is(err, identitydomain.ErrSessionLimit):
		return http.StatusTooManyRequests, "Limite de sessoes ativas atingido."
	default:
		return http.StatusBadRequest, "Nao foi possivel autenticar."
	}
}
