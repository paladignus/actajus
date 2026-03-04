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
	if !uc.checkSession {
		now := uc.clock.Now()
		active, err := uc.session.IsActive(ctx, claims.IDSession, claims.IDUser, now)
		if err != nil {
			return nil, err
		}
		if !active {
			return nil, domain.ErrSessionNotActive
		}
	}
	return &dto.AccessTokenClaims{
		IDUser:    claims.IDUser,
		IDSession: claims.IDSession,
	}, nil
}
