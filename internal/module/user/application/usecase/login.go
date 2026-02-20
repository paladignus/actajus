// Packge usecase
package usecase

import (
	"context"
	"time"

	sessionDomain "github.com/paladignus/actajus/internal/module/session/domain"
	"github.com/paladignus/actajus/internal/module/user/application/dto"
	"github.com/paladignus/actajus/internal/module/user/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
	"github.com/paladignus/actajus/internal/shared/infrastructure/security"
)

type Login struct {
	user    domain.UserRepository
	hash    domain.UserPasswordHashRepository
	session sessionDomain.SessionRepository

	MaxSessionsPerUser int
	RefreshTTL         time.Duration
	AccessTTL          time.Duration
}

func NewLogin(
	user domain.UserRepository,
	hash domain.UserPasswordHashRepository,
	session sessionDomain.SessionRepository,
) Login {
	return Login{
		user:    user,
		hash:    hash,
		session: session,
	}
}

func (uc Login) Execute(ctx context.Context, input dto.LoginRequest) (*dto.LoginReadModel, error) {
	now := time.Now()
	// ver a necessidade to tipo email ser um value object em vez de string
	user, err := uc.user.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, sharedDomain.ErrInvalidCredentials
	}
	if user.IsBlocked() {
		return nil, sharedDomain.ErrUserBlocked
	}
	if err := uc.hash.Compare(user.Password().Value(), input.Password); err != nil {
		return nil, sharedDomain.ErrInvalidCredentials
	}
	if uc.MaxSessionsPerUser > 0 {
		n, err := uc.session.CountActiveByUser(ctx, user.ID().Value(), now)
		if err != nil {
			return nil, err
		}
		if n >= uc.MaxSessionsPerUser {
			// Política: revogar as mais antigas / negar login.
			// Aqui: negar (simples). Você pode trocar por “revoke oldest”.
			return nil, sharedDomain.ErrSessionRevoked
		}
	}
	refresh, err := security.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshHash := security.HashToken(refresh)
	sid := vo.ID(9)
	refreshExp := now.Add(uc.RefreshTTL)
	s, err := sessionDomain.NewSessionBuilder().
		WithID(sid.Value()).
		WithIDUser(user.ID().Value()).
		WithRefreshTokenHash(refreshHash).
		WithExpiresAt(refreshExp).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		WithIP(input.IP).
		WithUserAgent(input.UserAgent).
		Build()
	if err != nil {
		return nil, err
	}
	if err := uc.session.Create(ctx, *s); err != nil {
		return nil, err
	}
	accessExp := now.Add(uc.AccessTTL)
	access, err := uc.Tokens.SignAccessToken(ports.AccessClaims{
		UserID:    string(u.ID),
		SessionID: string(sid),
		Roles:     u.Roles,
		ExpiresAt: accessExp,
		Issuer:    uc.Issuer,
		Audience:  uc.Audience,
	})
	if err != nil {
		return nil, err
	}

	return &dto.LoginReadModel{
		AccessToken:      access,
		RefreshToken:     refresh,
		IDSession:        string(sid),
		AccessExpiresAt:  accessExp.Unix(),
		RefreshExpiresAt: refreshExp.Unix(),
	}, nil
}
