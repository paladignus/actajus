// Package domain
package domain

import (
	"context"
	"time"
)

type SessionRepository interface {
	Create(ctx context.Context, session Session) error
	GetByID(ctx context.Context, sid string) (*Session, error)
	RotateRefreshToken(ctx context.Context, sid, newHash string, newExpiry, now time.Time) error
	Revoke(ctx context.Context, sid string, now time.Time) error
	RevokeAllByUser(ctx context.Context, idUser int64, now time.Time) error
	CountActiveByUser(ctx context.Context, idUser int64, now time.Time) (int, error)
}
