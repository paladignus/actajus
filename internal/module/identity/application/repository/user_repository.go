// Package repository
package repository

import (
	"context"

	identity "github.com/paladignus/actajus/internal/module/identity/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email vo.Email) (*identity.User, error)
	FindByID(ctx context.Context, id identity.IDUser) (*identity.User, error)
}
