// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Authenticate struct {
	persistencia repository.Account
	logger       repository.Logger
	// token        repository.Token
}

func NewAuthenticate(repository repository.Account, logger repository.Logger) Authenticate {
	return Authenticate{repository, logger}
}

func (a Authenticate) Execute(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error) {
	a.logger.Info(ctx, "authenticate user", "cpf", cpf)
	person, err := a.persistencia.FindPersonAccountByCPF(ctx, cpf)
	if err != nil {
		a.logger.Warn(ctx, "user not found during authentication",
			"cpf", cpf,
			"error", err,
		)
		return dto.AuthenticatedOutput{}, err
	}
	if ok := a.persistencia.ValidatePassword(ctx, person.IDUser, password); !ok {
		a.logger.Warn(ctx, "invalid credentials provided",
			"cpf", cpf,
			"id_person", person.IDUser,
		)
		return dto.AuthenticatedOutput{}, domain.ErrInvalidCredentials
	}
	a.logger.Info(ctx, "person authenticated successfully",
		"cpf", cpf,
		"id_person", person.IDUser,
	)
	return person, err

	// person, err := s.repository.Authenticate(ctx, cpf, password)
	// if err != nil {
	// 	return dto.AuthenticatedOutput{}, err
	// }
	// tokenPair, err := s.service.GenerateTokenPair(person.UserID)
	// if err != nil {
	// 	return dto.AuthenticatedOutput{}, err
	// }
	// person.AccessToken = tokenPair.AccessToken
	// person.RefreshToken = tokenPair.RefreshToken
	// return person, nil
}
