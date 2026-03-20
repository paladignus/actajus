// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
)

type IUser interface {
	AuthenticationByCPF(ctx context.Context, input command.SignInCommand) (user readmodel.SignInReadModel, err error)
	FindEmailByCPF(ctx context.Context, cpf string) (email readmodel.GetEmailByCPFReadModel, err error)
	FindIDUserByEmail(ctx context.Context, email string) (IDUser int, err error)
	UpdatePassword(ctx context.Context, idUser int, password, cpf string) error
}
