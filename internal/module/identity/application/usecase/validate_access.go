// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/application/mapper"
	"github.com/paladignus/actajus/internal/module/identity/application/repository"
	"github.com/paladignus/actajus/internal/module/identity/application/service"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type ValidateAccess struct {
	access       service.AccessTokenService
	session      repository.SessionRepository
	clock        service.Clock
	mapper       mapper.AuthMapper
	CheckSession bool
}

func NewValidateAccess(
	access service.AccessTokenService,
	session repository.SessionRepository,
	clock service.Clock,
	mapper mapper.AuthMapper,
	checkSession bool,
) ValidateAccess {
	return ValidateAccess{
		access: access, session: session, clock: clock,
		mapper: mapper, CheckSession: checkSession,
	}
}

func (uc ValidateAccess) Execute(ctx context.Context, input dto.ValidateAccessCommand) (*dto.AuthContext, error) {
	norm, err := uc.mapper.ValidateAccessInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid access token data: %w", err)
	}
	claims, err := uc.access.Verify(norm.Token)
	if err != nil {
		return nil, identity.ErrInvalidToken
	}
	now := uc.clock.Now()
	if !now.Before(claims.ExpiresAt) {
		return nil, identity.ErrInvalidToken
	}
	if uc.CheckSession {
		sid := identity.IDSession(claims.IDSession)
		sess, err := uc.session.GetByID(ctx, sid)
		if err != nil {
			return nil, identity.ErrSessionNotFound
		}
		if sess.IsRevoked() {
			return nil, identity.ErrSessionRevoked
		}
		if sess.IsExpired(now) {
			return nil, identity.ErrSessionExpired
		}
		if sess.IDUser().Value() != claims.IDUser {
			return nil, identity.ErrInvalidToken
		}
	}
	return &dto.AuthContext{
		IDUser:    claims.IDUser,
		IDSession: claims.IDSession,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}
