// Package repository
package repository

import (
	"context"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*identity.User, error)
	FindByID(ctx context.Context, id int64) (*identity.User, error)
	UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) error
}
