// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Account interface {
	FindUserAccountByCPF(context.Context, string) (dto.AuthenticateOutput, error)
	ValidatePassword(context.Context, string, string) error
	FindEmailByCPF(context.Context, string) (string, error)
}
