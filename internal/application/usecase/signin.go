// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type SignIn struct {
	user repository.IUser
	jwt  gateway.JWT
}

func NewSignIn(
	user repository.IUser,
	jwt gateway.JWT,
) SignIn {
	return SignIn{
		user,
		jwt,
	}
}

func (a SignIn) Execute(ctx context.Context, input dto.SignInInput) (user dto.SignInOutput, err error) {
	cpf := vo.CPF(input.CPF)
	if !cpf.IsValid() {
		return user, fmt.Errorf("invalid CPF format %s in sign in use case: %w", input.CPF, exception.ErrInvalidCredentials)
	}
	input.CPF = cpf.OnlyDigits()
	user, err = a.user.AuthenticationByCPF(ctx, input)
	if err != nil {
		return user, fmt.Errorf("sign in use case failed for CPF %s: %w", input.CPF, err)
	}
	tokenPair, err := a.jwt.GenerateTokenPair(string(user.IDUser))
	if err != nil {
		return user, fmt.Errorf("sign in use case failed to generate token pair for user ID %s: %w", user.IDUser, err)
	}
	user.AccessToken = tokenPair.AccessToken
	user.RefreshToken = tokenPair.RefreshToken
	return user, err
}
