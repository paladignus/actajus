// Package repository
package repository

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/module/identity/application/model"
	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, input *model.PasswordResetTokenCreate) (identity.IDPasswordReset, error)
	GetActiveByUser(ctx context.Context, userID identity.IDUser) (*model.PasswordResetToken, error)
	GetByID(ctx context.Context, id identity.IDPasswordReset) (*model.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id identity.IDPasswordReset, usedAt time.Time) error
	RevokeAllByUser(ctx context.Context, userID identity.IDUser, now time.Time) error
}
