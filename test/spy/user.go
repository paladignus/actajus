// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Result struct {
	User      dto.SignInOutput
	FindEmail dto.GetEmailByCPFOutput
}

type User struct {
	CallCount             int
	FindResult            Result
	FindError             error
	InvalidAllTokensError error
	ValidateError         error
}

func NewUser() *User {
	return &User{}
}

func (u *User) SignIn(ctx context.Context, cpf string) (authentication dto.SignInOutput, err error) {
	return u.FindResult.User, u.FindError
}

func (u *User) FindForAuthenticationByCPF(ctx context.Context, cpf string) (user dto.SignInOutput, err error) {
	return u.FindResult.User, u.FindError
}

// func (u *User) ValidatePassword(ctx context.Context, IDPeople, password string) (err error) {
// 	return u.ValidateError
// }

func (u *User) FindEmailByCPF(ctx context.Context, cpf string) (dto.GetEmailByCPFOutput, error) {
	return u.FindResult.FindEmail, u.FindError
}

func (u *User) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	return u.FindResult.User, u.FindError
}

func (u *User) FindByID(ctx context.Context, idAuthentication int) (entity.User, error) {
	return u.FindResult.User, u.FindError
}

func (u *User) UpdatePassword(ctx context.Context, idAuthentication int, hashedPassword string) error {
	return u.ValidateError
}

// func (u *User) AccountIsActive(context.Context, string) (string, error) {
// 	return u.FindResult.User.IDUser, u.FindError
// }
//
// func (u *User) InvalidAllTokensByIDUser(context.Context, string) error {
// 	return u.InvalidAllTokensError
// }
//
// func (u *User) CreateRecoverPassword(ctx context.Context, IDUser, token string) (err error) {
// 	return u.ValidateError
// }
