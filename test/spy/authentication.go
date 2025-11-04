// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type Result struct {
	Authentication dto.SignInOutput
	FindEmail      dto.GetEmailByCPFOutput
}

type AuthenticationSpy struct {
	CallCount     int
	FindResult    Result
	FindError     error
	InvalidError  error
	ValidateError error
}

func NewAuthenticationSpy() *AuthenticationSpy {
	return &AuthenticationSpy{}
}

func (a *AuthenticationSpy) SignIn(ctx context.Context, cpf string) (authentication dto.SignInOutput, err error) {
	return a.FindResult.Authentication, a.FindError
}

func (a *AuthenticationSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *AuthenticationSpy) FindEmailByCPF(ctx context.Context, cpf string) (dto.GetEmailByCPFOutput, error) {
	return a.FindResult.FindEmail, a.FindError
}

func (a *AuthenticationSpy) AccountIsActive(context.Context, string) (string, error) {
	return a.FindResult.Authentication.IDUser, a.FindError
}

func (a *AuthenticationSpy) InvalidAllTokensByIDUser(context.Context, string) error {
	return a.InvalidError
}

func (a *AuthenticationSpy) CreateRecoverPassword(ctx context.Context, IDUser, token string) (err error) {
	return a.ValidateError
}
