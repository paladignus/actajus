// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Authenticate struct {
	persistence repository.Account
	logger      repository.Logger
	gateway     gateway.Token
}

func NewAuthenticate(
	persistence repository.Account,
	logger repository.Logger,
	token gateway.Token,
) Authenticate {
	return Authenticate{
		persistence,
		logger,
		token,
	}
}

func (a Authenticate) Execute(ctx context.Context, req dto.AuthenticateInput) (dto.AuthenticateOutput, error) {
	a.logger.Info(ctx, "authenticate user", "cpf", req.CPF)
	cpf := vo.CPF(req.CPF)
	if !cpf.IsValid() {
		a.logger.Warn(ctx, "invalid cpf format provided",
			"cpf", req.CPF,
		)
		return dto.AuthenticateOutput{}, exception.ErrInvalidCPF
	}
	person, err := a.persistence.FindUserAccountByCPF(ctx, cpf.OnlyDigits())
	if err != nil {
		a.logger.Warn(ctx, "user not found during authentication",
			"cpf", req.CPF,
			"error", err,
		)
		return dto.AuthenticateOutput{}, err
	}
	err = a.persistence.ValidatePassword(ctx, person.IDUser, req.Password)
	if err != nil {
		a.logger.Warn(ctx, "invalid credentials provided",
			"cpf", req.CPF,
			"id_person", person.IDUser,
		)
		return dto.AuthenticateOutput{}, err
	}
	a.logger.Info(ctx, "person authenticated successfully",
		"cpf", req.CPF,
		"id_person", person.IDUser,
	)
	tokenPair, err := a.gateway.GenerateTokenPair(person.IDUser)
	if err != nil {
		a.logger.Error(ctx, "failed to generate token pair", "error", err)
	}
	person.AccessToken = tokenPair.AccessToken
	person.RefreshToken = tokenPair.RefreshToken
	a.logger.Info(ctx, "token pair generated successfully", "cpf", req.CPF)
	return person, err
}
