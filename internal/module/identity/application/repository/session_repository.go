// Package repository
package repository

import (
	"context"
	"time"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *identity.Session) error
	GetByID(ctx context.Context, sid identity.IDSession) (*identity.Session, error)
	RotateRefreshToken(ctx context.Context, sid identity.IDSession, newHash [32]byte, newExpiryAtTime time.Time) error
	Revoke(ctx context.Context, sid identity.IDSession) error
	RevokeAllByUser(ctx context.Context, uid identity.IDUser) error
	CountActiveByUser(ctx context.Context, uid identity.IDUser) (int, error)
	IsActive(ctx context.Context, sid identity.IDSession, uid identity.IDUser, now time.Time) (bool, error)
	RotateRefreshTokenAtomic(ctx context.Context, sid identity.IDSession, expectedOldHash [32]byte, newHash [32]byte, newExpiresAt time.Time, now time.Time) (rotated bool, err error)
}
