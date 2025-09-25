// Package spy
package spy

import (
	"context"
	"errors"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/domainerrors"
)

type AuthenticateSpy struct {
	ShouldReturnError        bool
	ShouldReturnUserNotFound bool
	CustomOutput             dto.AuthenticatedOutput
	CallCount                int
	LastCPF                  string
	LastPassword             string
}

func NewAuthenticateSpy() *AuthenticateSpy {
	return &AuthenticateSpy{
		CustomOutput: dto.AuthenticatedOutput{
			ID:        "1",
			FirstName: "John",
			LastName:  "Doe",
		},
	}
}

func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedOutput, error) {
	a.CallCount++
	a.LastCPF = cpf
	a.LastPassword = password
	if a.ShouldReturnUserNotFound {
		return dto.AuthenticatedOutput{}, domainerrors.ErrUserNotFound
	}
	if a.ShouldReturnError {
		return dto.AuthenticatedOutput{}, errors.New("spy database error")
	}
	return a.CustomOutput, nil
}
