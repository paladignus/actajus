// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
)

type AccountSpy struct {
	ShouldReturnError        bool
	ShouldReturnUserNotFound bool
	CustomOutput             dto.AuthenticateOutput
	CallCount                int
	LastCPF                  string
	LastPassword             string
}

func NewAccountSpy() *AccountSpy {
	return &AccountSpy{
		CustomOutput: dto.AuthenticateOutput{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			IDUser:       "1",
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "tes@email.com",
			Roles:        []string{"root"},
		},
	}
}

func (a AccountSpy) FindUserAccountByCPF(ctx context.Context, cpf string) (authenticated dto.AuthenticateOutput, err error) {
	return authenticated, err
}

func (a AccountSpy) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
	return err
}

// func (a *AuthenticateSpy) SignIn(ctx context.Context, cpf string, password string) (dto.AuthenticatedResponse, error) {
// 	a.CallCount++
// 	a.LastCPF = cpf
// 	a.LastPassword = password
// 	if a.ShouldReturnUserNotFound {
// 		return dto.AuthenticatedResponse{}, domain.ErrUserNotFound
// 	}
// 	if a.ShouldReturnError {
// 		return dto.AuthenticatedResponse{}, errors.New("spy database error")
// 	}
// 	return a.CustomOutput, nil
// }
