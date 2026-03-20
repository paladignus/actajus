// Package service provides domain services for identity module
package service

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
)

// RefreshTokenGenerator generates refresh tokens
type RefreshTokenGenerator interface {
	Generate() (token string, hash [32]byte, err error)
}

// AccessTokenSigner signs access tokens
type AccessTokenSigner interface {
	Sign(claims dto.AccessTokenClaims) (string, error)
}

// UserRepository provides user persistence
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

// SessionRepository provides session persistence
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	CountActiveByUser(ctx context.Context, idUser int64) (int, error)
}

// AuthnService handles authentication logic
type AuthnService struct {
	user    UserRepository
	session SessionRepository
	hasher  PasswordHasher
	refresh RefreshTokenGenerator
	access  AccessTokenSigner
	config  config.AuthConfig
	clock   domain.Clock
}

// NewAuthnService creates a new authentication service
func NewAuthnService(
	user UserRepository,
	session SessionRepository,
	hasher PasswordHasher,
	refresh RefreshTokenGenerator,
	access AccessTokenSigner,
	config config.AuthConfig,
	clock domain.Clock,
) AuthnService {
	return AuthnService{user, session, hasher, refresh, access, config, clock}
}

// AuthenticateResult holds authentication result
type AuthenticateResult struct {
	User         *domain.User
	Session      *domain.Session
	AccessToken  string
	RefreshToken string
	AccessExp    time.Time
	RefreshExp   time.Time
}

// Authenticate performs user authentication
func (s AuthnService) Authenticate(
	ctx context.Context,
	email string,
	password string,
	ip string,
	userAgent string,
) (*AuthenticateResult, error) {
	// Find user
	user, err := s.user.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is blocked
	if user.IsBlocked() {
		return nil, domain.ErrUserBlocked
	}

	// Verify password
	if !user.MatchesPassword(password, s.hasher) {
		return nil, domain.ErrInvalidCredentials
	}

	// Check session limit
	if s.config.MaxSessions > 0 {
		n, err := s.session.CountActiveByUser(ctx, user.ID().Value())
		if err != nil {
			return nil, err
		}
		if n >= s.config.MaxSessions {
			return nil, domain.ErrSessionLimit
		}
	}

	// Generate tokens
	refreshToken, refreshHash, err := s.refresh.Generate()
	if err != nil {
		return nil, err
	}

	now := s.clock.Now()
	refreshExp := now.Add(s.config.RefreshTTL)

	// Create session
	session, err := domain.NewSessionBuilder().
		WithIDUser(user.ID().Value()).
		WithRefreshHash(refreshHash).
		WithExpiresAt(refreshExp).
		WithIP(ip).
		WithUserAgent(userAgent).
		Build()
	if err != nil {
		return nil, err
	}

	if err := s.session.Create(ctx, session); err != nil {
		return nil, err
	}

	// Sign access token
	accessExp := now.Add(s.config.AccessTTL)
	accessToken, err := s.access.Sign(dto.AccessTokenClaims{
		IDSession: session.ID().Value(),
		IDUser:    user.ID().Value(),
		Issuer:    s.config.Issuer,
		Audience:  s.config.Audience,
		IssuedAt:  now,
		ExpiresAt: accessExp,
	})
	if err != nil {
		return nil, err
	}

	return &AuthenticateResult{
		User:         user,
		Session:      session,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExp:    accessExp,
		RefreshExp:   refreshExp,
	}, nil
}
