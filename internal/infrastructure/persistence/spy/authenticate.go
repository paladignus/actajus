// Package spy
package spy

import (
	"context"
	"errors"

	"github.com/paladignus/actajus/internal/application/dto"
)

type AuthenticateSpy struct {
	CPF        string
	Password   string
	CallsCount int
}

func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf, password string) (dto.AuthenticatedOutput, error) {
	authenticated := dto.AuthenticatedOutput{}
	if cpf == "" || password == "" {
		return authenticated, errors.New("cpf and password cannot be empty")
	}
	a.CPF = cpf
	a.Password = password
	a.CallsCount++
	authenticated.ID = "1"
	authenticated.FirstName = "John"
	authenticated.LastName = "Doe"
	return authenticated, nil
}
