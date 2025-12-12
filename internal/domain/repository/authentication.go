// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Authentication interface {
	SignIn(context.Context, string) (dto.SignInOutput, error)
	ValidatePassword(context.Context, string, string) error
	FindEmailByCPF(context.Context, string) (dto.GetEmailByCPFOutput, error)
	AccountIsActive(context.Context, string) (string, error)
	InvalidAllTokensByIDUser(context.Context, string) error
	CreateRecoverPassword(context.Context, string, string) error

	FindByEmail(ctx context.Context, email string) (entity.Authentication, error)
	FindById(ctx context.Context, idAuthentication int) (entity.Authentication, error)
	UpdatePassword(ctx context.Context, idAuthentication int, hashedPassword string) error
}
