// Package adapter
package adapter

import (
	"errors"

	"connectrpc.com/connect"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

func ToConnectIdentityError(err error) *connect.Error {
	if err == nil {
		return nil
	}
	var ve sharedDomain.ValidationError
	if errors.As(err, &ve) {
		return ValidationErrorDetail(ve)
	}
	switch {
	case errors.Is(err, identity.ErrInvalidCredentials):
		return connect.NewError(connect.CodeUnauthenticated, err)
	case errors.Is(err, identity.ErrInvalidToken):
		return connect.NewError(connect.CodeUnauthenticated, err)
	case errors.Is(err, identity.ErrSessionExpired),
		errors.Is(err, identity.ErrSessionRevoked),
		errors.Is(err, identity.ErrSessionNotFound):
		return connect.NewError(connect.CodeUnauthenticated, err)
	case errors.Is(err, identity.ErrUserBlocked):
		return connect.NewError(connect.CodePermissionDenied, err)
	case errors.Is(err, identity.ErrSessionLimit):
		return connect.NewError(connect.CodeResourceExhausted, err)
	case errors.Is(err, identity.ErrRefreshReuse):
		return connect.NewError(connect.CodeUnauthenticated, err)
	default:
		// Depois mapear melhor pra detectar pgx/redis
		return connect.NewError(connect.CodeInternal, err)
	}
}
