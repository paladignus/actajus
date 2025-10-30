// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Authentication interface {
	SignIn(context.Context, string) (dto.SignInOutput, error)
	ValidatePassword(context.Context, string, string) error
	FindEmailByCPF(context.Context, string) (dto.GetEmailByCPFOutput, error)
}
