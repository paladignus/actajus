// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/domain/service"
)

type SignIn struct {
	repository repository.Authenticate
	service    service.Token
}

func NewSignIn(repository repository.Authenticate, service service.Token) SignIn {
	return SignIn{repository, service}
}

func (s SignIn) Execute(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error) {
	person, err := s.repository.SignIn(ctx, cpf, password)
	if err != nil {
		return dto.AuthenticatedOutput{}, err
	}
	tokenPair, err := s.service.GenerateTokenPair(person.UserID)
	if err != nil {
		return dto.AuthenticatedOutput{}, err
	}
	person.AccessToken = tokenPair.AccessToken
	person.RefreshToken = tokenPair.RefreshToken
	return person, nil
}
