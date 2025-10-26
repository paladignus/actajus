// Package usecase
package usecase

import (
	"context"

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

func (g GetEmailByCPF) Execute(ctx context.Context, cpfString string) (string, error) {
	g.logger.Info(ctx, "getting email by CPF", "cpf", cpfString)
	cpf := vo.CPF(cpfString)
	if !cpf.IsValid() {
		g.logger.Warn(ctx, "invalid cpf format provided", "cpf", cpfString)
		return "", exception.ErrInvalidCPF
	}
	email, err := g.persistence.GetEmailByCPF(ctx, cpf.OnlyDigits())
	if err != nil {
		g.logger.Warn(ctx, "email not found for provided cpf", "cpf", cpfString, "error", err)
		return "", exception.ErrEmailNotFound
	}
	g.logger.Info(ctx, "email found successfully", "cpf", cpfString, "email", email)
	return email, nil
}
