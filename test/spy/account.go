// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type AccountSpy struct {
	CallCount     int
	FindResult    dto.AuthenticateOutput
	FindError     error
	ValidateError error
}

func NewAccountSpy() *AccountSpy {
	return &AccountSpy{}
}

func (a *AccountSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return a.FindResult, a.FindError
}

func (a *AccountSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *AccountSpy) FindEmailByCPF(ctx context.Context, cpf string) (string, error) {
	return a.FindResult.Email, a.FindError
}
