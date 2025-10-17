// Package controller
package controller

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/application/usecase"
)

type Authenticate struct {
	usecase usecase.Authenticate
}

func NewAuthenticate(usecase usecase.Authenticate) Authenticate {
	return Authenticate{usecase}
}

func (a Authenticate) Authenticate(ctx context.Context, input dto.AuthenticateInput) (dto.AuthenticateOutput, error) {
	return a.usecase.Execute(ctx, input)
}
