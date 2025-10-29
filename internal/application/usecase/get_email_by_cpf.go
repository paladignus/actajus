// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type GetEmailByCPF struct {
	persistence repository.Account
	logger      repository.Logger
}

func NewGetEmailByCPF(persistence repository.Account, logger repository.Logger) GetEmailByCPF {
	return GetEmailByCPF{persistence, logger}
}

func (g GetEmailByCPF) Execute(ctx context.Context, req dto.GetEmailByCPFInput) (dto.GetEmailByCPFOutput, error) {
	g.logger.Info(ctx, "getting email by CPF", "cpf", req.CPF)
	cpf := vo.CPF(req.CPF)
	if !cpf.IsValid() {
		g.logger.Warn(ctx, "invalid cpf format provided", "cpf", req.CPF)
		return dto.GetEmailByCPFOutput{}, exception.ErrInvalidCPF
	}
	email, err := g.persistence.FindEmailByCPF(ctx, cpf.OnlyDigits())
	if err != nil {
		g.logger.Warn(ctx, "email not found for provided cpf", "cpf", req.CPF, "error", err)
		return dto.GetEmailByCPFOutput{}, exception.ErrEmailNotFound
	}
	g.logger.Info(ctx, "email found successfully", "cpf", req.CPF, "email", email)
	return dto.GetEmailByCPFOutput{Email: email}, nil
}
