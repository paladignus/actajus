// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Session struct {
	id               vo.ID
	idUser           vo.ID
	refreshTokenHash string
	userAgent        string
	ip               string
	expiresAt        time.Time
	createdAt        time.Time
	updatedAt        time.Time
	revokedAt        *time.Time
}

type SessionBuilder struct {
	session *Session
}

func NewSessionBuilder() *SessionBuilder {
	now := time.Now()
	return &SessionBuilder{
		session: &Session{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (s *SessionBuilder) WithID(id int64) *SessionBuilder {
	s.session.id = vo.ID(id)
	return s
}

func (s *SessionBuilder) WithIDUser(id int64) *SessionBuilder {
	s.session.idUser = vo.ID(id)
	return s
}

func (s *SessionBuilder) WithRefreshTokenHash(hash string) *SessionBuilder {
	s.session.refreshTokenHash = hash
	return s
}

func (s *SessionBuilder) WithUserAgent(ua string) *SessionBuilder {
	s.session.userAgent = ua
	return s
}

func (s *SessionBuilder) WithIP(ip string) *SessionBuilder {
	s.session.ip = ip
	return s
}

func (s *SessionBuilder) WithExpiresAt(expiry time.Time) *SessionBuilder {
	s.session.expiresAt = expiry
	return s
}

func (s *SessionBuilder) WithCreatedAt(createdAt time.Time) *SessionBuilder {
	s.session.createdAt = createdAt
	return s
}

func (s *SessionBuilder) WithUpdatedAt(updatedAt time.Time) *SessionBuilder {
	s.session.updatedAt = updatedAt
	return s
}

func (s *SessionBuilder) WithRevokedAt(revokedAt *time.Time) *SessionBuilder {
	s.session.revokedAt = revokedAt
	return s
}

func (s *SessionBuilder) Build() (*Session, error) {
	return s.session, nil
}

func (s *Session) ID() vo.ID                { return s.id }
func (s *Session) IDUser() vo.ID            { return s.idUser }
func (s *Session) RefreshTokenHash() string { return s.refreshTokenHash }
func (s *Session) UserAgent() string        { return s.userAgent }
func (s *Session) IP() string               { return s.ip }
func (s *Session) ExpiresAt() time.Time     { return s.expiresAt }
func (s *Session) CreatedAt() time.Time     { return s.createdAt }
func (s *Session) UpdatedAt() time.Time     { return s.updatedAt }
func (s *Session) RevokedAt() *time.Time    { return s.revokedAt }

// func (s *Session) isExpired() bool { return s.expiresAt.Before(time.Now()) }

func (s *Session) IsRevoked() bool              { return s.revokedAt != nil }
func (s *Session) IsExpired(now time.Time) bool { return !now.Before(s.expiresAt) }
func (s *Session) IsActive(now time.Time) bool  { return !s.IsRevoked() && !s.IsExpired(now) }
