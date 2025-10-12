// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Account interface {
	FindUserAccountByCPF(context.Context, string) (dto.AuthenticatedResponse, error)
	ValidatePassword(context.Context, string, string) error
}
