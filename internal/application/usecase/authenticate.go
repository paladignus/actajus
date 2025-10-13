// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Authenticate struct {
	repository.Account
	repository.Logger
	gateway.Token
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
	a.Info(ctx, "authenticate user", "cpf", req.CPF)
	person, err := a.FindUserAccountByCPF(ctx, req.CPF)
	if err != nil {
		a.Warn(ctx, "user not found during authentication",
			"cpf", req.CPF,
			"error", err,
		)
		return dto.AuthenticateOutput{}, err
	}
	err = a.ValidatePassword(ctx, person.IDUser, req.Password)
	if err != nil {
		a.Warn(ctx, "invalid credentials provided",
			"cpf", req.CPF,
			"id_person", person.IDUser,
		)
		return dto.AuthenticateOutput{}, err
	}
	a.Info(ctx, "person authenticated successfully",
		"cpf", req.CPF,
		"id_person", person.IDUser,
	)
	tokenPair, err := a.GenerateTokenPair(person.IDUser)
	if err != nil {
		a.Error(ctx, "failed to generate token pair", "error", err)
	}
	person.AccessToken = tokenPair.AccessToken
	person.RefreshToken = tokenPair.RefreshToken
	a.Info(ctx, "token pair generated successfully", "cpf", req.CPF)
	return person, err
}
