// Package repository
package repository

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/model"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, input *model.PasswordResetTokenCreate) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id int64, usedAt time.Time) error
	RevokeAllByUser(ctx context.Context, uid int64, now time.Time) error
	GetActiveByUser(ctx context.Context, uid int64) (*model.PasswordResetToken, error)
}
