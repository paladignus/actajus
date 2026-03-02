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
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
)

type Refresh struct {
	session    repository.SessionRepository
	user       repository.UserRepository
	refresh    service.RefreshTokenService
	access     service.AccessTokenService
	clock      service.Clock
	config     config.AuthConfig
	mapper     mapper.AuthMapper
	projection mapper.AuthProjectionMapper
}

func NewRefresh(
	session repository.SessionRepository,
	user repository.UserRepository,
	refresh service.RefreshTokenService,
	access service.AccessTokenService,
	clock service.Clock,
	config config.AuthConfig,
	mapper mapper.AuthMapper,
	projection mapper.AuthProjectionMapper,
) Refresh {
	return Refresh{
		session, user, refresh, access, clock,
		config, mapper, projection,
	}
}

func (uc Refresh) Execute(ctx context.Context, input dto.RefreshCommand) (*dto.AuthTokensReadModel, error) {
	norm, err := uc.mapper.RefreshInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh data: %w", err)
	}
	now := uc.clock.Now()
	sid := identity.IDSession(norm.IDSession)
	sess, err := uc.session.GetByID(ctx, sid)
	if err != nil {
		return nil, identity.ErrSessionNotFound
	}
	if sess.IsRevoked() {
		return nil, identity.ErrSessionRevoked
	}
	if sess.IsExpired(now) {
		_ = uc.session.Revoke(ctx, sid)
		return nil, identity.ErrSessionExpired
	}
	oldHash, ok := uc.refresh.Hash(norm.RefreshToken)
	if !ok {
		_ = uc.session.Revoke(ctx, sid)
		return nil, identity.ErrRefreshReuse
	}
	newRefreshToken, newHash, err := uc.refresh.Generate()
	if err != nil {
		return nil, err
	}
	newRefreshExp := now.Add(uc.config.RefreshTTL)
	rotated, err := uc.session.RotateRefreshTokenAtomic(ctx,
		sid,
		oldHash,
		newHash,
		newRefreshExp,
		now)
	if err != nil {
		return nil, err
	}
	if !rotated {
		_ = uc.session.Revoke(ctx, sid)
		return nil, identity.ErrRefreshReuse
	}
	accessExp := now.Add(uc.config.AccessTTL)
	accessToken, err := uc.access.Sign(dto.AccessTokenClaims{
		IDUser:    sess.IDUser().Value(),
		IDSession: sess.ID().Value(),
		Issuer:    uc.config.Issuer,
		Audience:  uc.config.Audience,
		IssuedAt:  now,
		ExpiresAt: accessExp,
	})
	if err != nil {
		return nil, err
	}
	user, err := uc.user.FindByID(ctx, sess.IDUser())
	if err != nil {
		_ = uc.session.Revoke(ctx, sid)
		return nil, identity.ErrInvalidToken
	}
	if user.IsBlocked() {
		_ = uc.session.Revoke(ctx, sid)
		return nil, identity.ErrUserBlocked
	}
	return uc.projection.ProjectTokens(
		sess.ID(),
		sess.IDUser(),
		accessToken,
		newRefreshToken,
		accessExp,
		newRefreshExp,
	), nil
}
