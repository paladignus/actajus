// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
)

type Session struct {
	id          IDSession
	idUser      IDUser
	refreshHash [32]byte
	expiresAt   time.Time
	revokedAt   *time.Time
	rotatedAt   *time.Time
	ip          string
	userAgent   string
	createdAt   time.Time
	updatedAt   time.Time
}

type SessionBuilder struct {
	s *Session
}

func NewSessionBuilder() *SessionBuilder {
	now := time.Now()
	return &SessionBuilder{
		&Session{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (b *SessionBuilder) WithID(id IDSession) *SessionBuilder {
	b.s.id = id
	return b
}

func (b *SessionBuilder) WithIDUser(id IDUser) *SessionBuilder {
	b.s.idUser = id
	return b
}

func (b *SessionBuilder) WithRefreshHash(h [32]byte) *SessionBuilder {
	b.s.refreshHash = h
	return b
}

func (b *SessionBuilder) WithExpiresAt(t time.Time) *SessionBuilder {
	b.s.expiresAt = t
	return b
}

func (b *SessionBuilder) WithRevokedAt(t *time.Time) *SessionBuilder {
	b.s.revokedAt = t
	return b
}

func (b *SessionBuilder) WithRotatedAt(t *time.Time) *SessionBuilder {
	b.s.rotatedAt = t
	return b
}

func (b *SessionBuilder) WithIP(ip string) *SessionBuilder {
	b.s.ip = ip
	return b
}

func (b *SessionBuilder) WithUserAgent(ua string) *SessionBuilder {
	b.s.userAgent = ua
	return b
}

func (b *SessionBuilder) WithCreatedAt(t time.Time) *SessionBuilder {
	b.s.createdAt = t
	return b
}

func (b *SessionBuilder) WithUpdatedAt(t time.Time) *SessionBuilder {
	b.s.updatedAt = t
	return b
}

func (b *SessionBuilder) Build() (*Session, error) {
	if b.s.idUser == 0 {
		return nil, domain.NewFieldError("user_id", "user id is invalid")
	}
	if b.s.expiresAt.IsZero() {
		return nil, domain.NewFieldError("expires_at", "expires_at is required")
	}
	return b.s, nil
}

func (s *Session) ID() IDSession                { return s.id }
func (s *Session) IDUser() IDUser               { return s.idUser }
func (s *Session) RefreshHash() [32]byte        { return s.refreshHash }
func (s *Session) ExpiresAt() time.Time         { return s.expiresAt }
func (s *Session) RevokedAt() *time.Time        { return s.revokedAt }
func (s *Session) RotatedAt() *time.Time        { return s.rotatedAt }
func (s *Session) IP() string                   { return s.ip }
func (s *Session) UserAgent() string            { return s.userAgent }
func (s *Session) IsRevoked() bool              { return s.revokedAt != nil }
func (s *Session) IsExpired(now time.Time) bool { return !now.Before(s.expiresAt) }
func (s *Session) CreatedAt() time.Time         { return s.createdAt }
func (s *Session) UpdatedAt() time.Time         { return s.updatedAt }

func (s *Session) SetID(id IDSession) error {
	if s.id != 0 {
		return domain.NewFieldError("id", "session id is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	s.id = id
	return nil
}

func (s *Session) Revoke(now time.Time) error {
	if s.id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	if s.IsRevoked() {
		return domain.NewFieldError("revoked_at", "session already revoked")
	}
	s.revokedAt = &now
	s.updatedAt = now
	return nil
}

func (s *Session) Rotate(newHash [32]byte, newExp time.Time, now time.Time) error {
	if s.id == 0 {
		return domain.NewFieldError("id", "session id is invalid")
	}
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	s.refreshHash = newHash
	s.expiresAt = newExp
	s.rotatedAt = &now
	s.updatedAt = now
	return nil
}
