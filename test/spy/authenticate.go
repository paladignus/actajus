// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type AuthenticateSpy struct {
	CallCount     int
	FindResult    dto.AuthenticateOutput
	FindError     error
	ValidateError error
}

func NewAuthenticateSpy() *AuthenticateSpy {
	return &AuthenticateSpy{}
}

func (a *AuthenticateSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return a.FindResult, a.FindError
}

func (a *AuthenticateSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *AuthenticateSpy) GetEmailByCPF(ctx context.Context, cpf string) (string, error) {
	return a.FindResult.Email, a.FindError
}
