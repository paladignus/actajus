// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type GetEmailByCPF struct {
	user repository.IUser
}

func NewGetEmailByCPF(
	persistence repository.IUser,
) GetEmailByCPF {
	return GetEmailByCPF{
		persistence,
	}
}

func (g GetEmailByCPF) Execute(
	ctx context.Context,
	input dto.GetEmailByCPFInput,
) (email dto.GetEmailByCPFOutput, err error) {
	cpf := vo.CPF(input.CPF)
	if !cpf.IsValid() {
		return email, fmt.Errorf("use case get email by cpf, invalid cpf: %w", exception.ErrEmailNotFound)
	}
	email, err = g.user.FindEmailByCPF(ctx, cpf.OnlyDigits())
	if err != nil {
		return email, fmt.Errorf("use case get email by cpf: %w", err)
	}
	return email, nil
}
