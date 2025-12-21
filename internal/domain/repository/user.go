// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type IUser interface {
	AuthenticationByCPF(ctx context.Context, input dto.SignInInput) (user dto.SignInOutput, err error)
	FindEmailByCPF(ctx context.Context, cpf string) (email dto.GetEmailByCPFOutput, err error)
	FindIDUserByEmail(ctx context.Context, email string) (IDUser int, err error)
	UpdatePassword(ctx context.Context, idUser int, password string) error
}
