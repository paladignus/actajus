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
	RotateRefresh(ctx context.Context, sid identity.IDSession, newHash [32]byte, newExpiryUnix time.Time) error
	Revoke(ctx context.Context, sid identity.IDSession) error
	RevokeAllByUser(ctx context.Context, uid identity.IDUser) error
	CountActiveByUser(ctx context.Context, uid identity.IDUser) (int, error)
}
