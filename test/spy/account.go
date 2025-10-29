// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Result struct {
	Authenticated   dto.AuthenticateOutput
	FindEmail       dto.GetEmailByCPFOutput
	RecoverPassword dto.RecoverPasswordOutput
}

type AccountSpy struct {
	CallCount     int
	FindResult    Result
	FindError     error
	ValidateError error
}

func NewAccountSpy() *AccountSpy {
	return &AccountSpy{}
}

func (a *AccountSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return a.FindResult.Authenticated, a.FindError
}

func (a *AccountSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *AccountSpy) FindEmailByCPF(ctx context.Context, cpf string) (dto.GetEmailByCPFOutput, error) {
	return a.FindResult.FindEmail, a.FindError
}
