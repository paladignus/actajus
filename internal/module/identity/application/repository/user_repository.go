// Package repository
package repository

import (
	"context"
	"time"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*identity.User, error)
	FindByID(ctx context.Context, id int64) (*identity.User, error)
	FindPrimaryEmailIDByUser(ctx context.Context, id int64) (int64, error)
	Create(ctx context.Context, input CreateUserRegistration) (CreatedUserRegistration, error)
	UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) error
	UpdateLastLoginAt(ctx context.Context, id int64, lastLoginAt time.Time) error
	MarkPrimaryEmailVerified(ctx context.Context, idEmail int64, verifiedAt time.Time) error
}

type CreateUserRegistration struct {
	FirstName    string
	LastName     string
	Birthday     string
	GenderID     uint
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreatedUserRegistration struct {
	IDUser  int64
	IDEmail int64
}
