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

type Authentication struct {
	CallCount             int
	FindResult            Result
	FindError             error
	InvalidAllTokensError error
	ValidateError         error
}

func NewAuthentication() *Authentication {
	return &Authentication{}
}

func (a *Authentication) SignIn(ctx context.Context, cpf string) (authentication dto.SignInOutput, err error) {
	return a.FindResult.Authentication, a.FindError
}

func (a *Authentication) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return a.ValidateError
}

func (a *Authentication) FindEmailByCPF(ctx context.Context, cpf string) (dto.GetEmailByCPFOutput, error) {
	return a.FindResult.FindEmail, a.FindError
}

func (a *Authentication) AccountIsActive(context.Context, string) (string, error) {
	return a.FindResult.Authentication.IDUser, a.FindError
}

func (a *Authentication) InvalidAllTokensByIDUser(context.Context, string) error {
	return a.InvalidAllTokensError
}

func (a *Authentication) CreateRecoverPassword(ctx context.Context, IDUser, token string) (err error) {
	return a.ValidateError
}
