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

type Login struct {
	user       repository.UserRepository
	session    repository.SessionRepository
	hasher     service.PasswordHasher
	refresh    service.RefreshTokenService
	access     service.AccessTokenService
	clock      service.Clock
	config     config.AuthConfig
	mapper     mapper.AuthMapper
	projection mapper.AuthProjectionMapper
}

func NewLogin(
	user repository.UserRepository,
	session repository.SessionRepository,
	hasher service.PasswordHasher,
	refresh service.RefreshTokenService,
	access service.AccessTokenService,
	clock service.Clock,
	config config.AuthConfig,
	mapper mapper.AuthMapper,
	projection mapper.AuthProjectionMapper,
) Login {
	return Login{
		user, session, hasher, refresh, access,
		clock, config, mapper, projection,
	}
}

func (uc Login) Execute(ctx context.Context, input dto.LoginCommand) (*dto.AuthTokensReadModel, error) {
	norm, err := uc.mapper.LoginInputToNormalized(input)
	if err != nil {
		return nil, fmt.Errorf("invalid login data: %w", err)
	}
	user, err := uc.user.FindByEmail(ctx, norm.Email)
	if err != nil {
		fmt.Println(err)
		return nil, identity.ErrInvalidCredentials
	}
	if user.IsBlocked() {
		return nil, identity.ErrUserBlocked
	}
	if err := uc.hasher.Compare(user.PasswordHash().Value(), norm.Password); err != nil {
		fmt.Println(err)
		return nil, identity.ErrInvalidCredentials
	}
	// if hasher.NeedsRehash(user.PasswordHash().Value()) {
	// 	newHash := hasher.Hash(plain);
	// 	repo.UpdatePasswordHash(...)
	// }
	// Fazer rehash detection no futuro
	if uc.config.MaxSessions > 0 {
		n, err := uc.session.CountActiveByUser(ctx, user.ID())
		if err != nil {
			return nil, err
		}
		if n >= uc.config.MaxSessions {
			return nil, identity.ErrSessionLimit
		}
	}
	refreshToken, refreshHash, err := uc.refresh.Generate()
	if err != nil {
		return nil, err
	}
	now := uc.clock.Now()
	refreshExp := now.Add(uc.config.RefreshTTL)
	sess, err := identity.NewSessionBuilder().
		WithIDUser(user.ID()).
		WithRefreshHash(refreshHash).
		WithExpiresAt(refreshExp).
		WithIP(norm.IP).
		WithUserAgent(norm.UserAgent).
		Build()
	if err != nil {
		return nil, err
	}
	if err := uc.session.Create(ctx, sess); err != nil {
		return nil, err
	}
	accessExp := now.Add(uc.config.AccessTTL)
	accessToken, err := uc.access.Sign(dto.AccessTokenClaims{
		IDSession: sess.ID().Value(),
		IDUser:    user.ID().Value(),
		Issuer:    uc.config.Issuer,
		Audience:  uc.config.Audience,
		IssuedAt:  now,
		ExpiresAt: accessExp,
	})
	if err != nil {
		return nil, err
	}
	return uc.projection.ProjectTokens(
		sess.ID(),
		user.ID(),
		accessToken,
		refreshToken,
		accessExp,
		refreshExp,
	), nil
}
