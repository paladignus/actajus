// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/application/readmodel"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Result struct {
	UserDTO    readmodel.SignInReadModel
	UserEntity entity.User
	FindEmail  readmodel.GetEmailByCPFReadModel
}

type User struct {
	CallCount     int
	FindResult    Result
	FindError     error
	ValidateError error
}

func NewUser() *User {
	return &User{}
}

func (u *User) AuthenticationByCPF(ctx context.Context, input command.SignInCommand) (user readmodel.SignInReadModel, err error) {
	return u.FindResult.UserDTO, u.FindError
}

func (u *User) FindEmailByCPF(ctx context.Context, cpf string) (readmodel.GetEmailByCPFReadModel, error) {
	return u.FindResult.FindEmail, u.FindError
}

func (u *User) FindIDUserByEmail(ctx context.Context, email string) (int, error) {
	return 0, u.FindError
	// return u.FindResult.UserEntity, u.FindError
}

// func (u *User) FindByID(ctx context.Context, idAuthentication int) (entity.User, error) {
// 	return u.FindResult.UserEntity, u.FindError
// }

func (u *User) UpdatePassword(ctx context.Context, idAuthentication int, hashedPassword, cpf string) error {
	return u.ValidateError
}
