// Package repository
package repository

import (
	"context"
	"time"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *identity.Session) error
	GetByID(ctx context.Context, id int64) (*identity.Session, error)
	RotateRefreshToken(ctx context.Context, id int64, hash [32]byte, expiryAtTime time.Time) error
	Revoke(ctx context.Context, id int64) error
	RevokeAllByUser(ctx context.Context, uid int64) error
	CountActiveByUser(ctx context.Context, uid int64) (int, error)
	IsActive(ctx context.Context, id int64, uid int64, now time.Time) (bool, error)
	RotateRefreshTokenAtomic(ctx context.Context, id int64, oldHash [32]byte, hash [32]byte, expiresAt time.Time, now time.Time) (rotated bool, err error)
}
