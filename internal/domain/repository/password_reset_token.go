// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type IPasswordResetToken interface {
	Create(ctx context.Context, token entity.PasswordResetToken) error
	FindByToken(ctx context.Context, token string) (entity.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, token string) error
	InvalidateUserTokens(ctx context.Context, idUser int) error
}
