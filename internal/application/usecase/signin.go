// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type SignIn struct {
	repository repository.Authenticate
}

func NewSignIn(repository repository.Authenticate) SignIn {
	return SignIn{repository}
}

func (s SignIn) Execute(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error) {
	return s.repository.SignIn(ctx, cpf, password)
}
