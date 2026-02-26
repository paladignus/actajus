// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	"github.com/paladignus/actajus/internal/module/identity/domain"
)

type ValidateAccess struct {
	token        service.AccessTokenService
	session      repository.SessionRepository
	clock        service.Clock
	checkSession bool
}

func NewValidateAccess(
	token service.AccessTokenService,
	session repository.SessionRepository,
	clock service.Clock,
	checkSession bool,
) ValidateAccess {
	return ValidateAccess{
		token, session,
		clock, checkSession,
	}
}

func (uc ValidateAccess) Execute(ctx context.Context, cmd dto.ValidateAccessCommand) (*dto.AccessTokenClaims, error) {
	if cmd.AccessToken == "" {
		return nil, domain.ErrMissingAccessToken
	}
	claims, err := uc.token.Verify(cmd.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err)
	}
	idUser := domain.IDUser(claims.IDUser)
	idSession := domain.IDSession(claims.IDSession)
	if !uc.checkSession {
		now := uc.clock.Now()
		active, err := uc.session.IsActive(ctx, idSession, idUser, now)
		if err != nil {
			return nil, err
		}
		if !active {
			return nil, domain.ErrSessionNotActive
		}
	}
	return &dto.AccessTokenClaims{
		IDUser:    idUser.Value(),
		IDSession: idSession.Value(),
	}, nil
}

// func (uc ValidateAccess) Execute(ctx context.Context, input dto.ValidateAccessCommand) (*dto.AuthContext, error) {
// 	norm, err := uc.mapper.ValidateAccessInputToNormalized(input)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid access token data: %w", err)
// 	}
// 	claims, err := uc.token.Verify(norm.Token)
// 	if err != nil {
// 		return nil, identity.ErrInvalidToken
// 	}
// 	now := uc.clock.Now()
// 	if !now.Before(claims.ExpiresAt) {
// 		return nil, identity.ErrInvalidToken
// 	}
// 	if uc.CheckSession {
// 		sid := identity.IDSession(claims.IDSession)
// 		sess, err := uc.session.GetByID(ctx, sid)
// 		if err != nil {
// 			return nil, identity.ErrSessionNotFound
// 		}
// 		if sess.IsRevoked() {
// 			return nil, identity.ErrSessionRevoked
// 		}
// 		if sess.IsExpired(now) {
// 			return nil, identity.ErrSessionExpired
// 		}
// 		if sess.IDUser().Value() != claims.IDUser {
// 			return nil, identity.ErrInvalidToken
// 		}
// 	}
// 	return &dto.AuthContext{
// 		IDUser:    claims.IDUser,
// 		IDSession: claims.IDSession,
// 		IssuedAt:  claims.IssuedAt,
// 		ExpiresAt: claims.ExpiresAt,
// 	}, nil
// }
