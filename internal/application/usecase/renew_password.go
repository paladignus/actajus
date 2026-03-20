// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type RenewPassword struct {
	user  repository.IUser
	token repository.IPasswordResetToken
}

func NewRenewPassword(
	user repository.IUser,
	token repository.IPasswordResetToken,
) RenewPassword {
	return RenewPassword{
		user,
		token,
	}
}

func (r RenewPassword) Execute(ctx context.Context, input command.RenewPasswordCommand) error {
	cpf := vo.CPF(input.CPF)
	if !cpf.IsValid() {
		return fmt.Errorf("use case renew password, invalid cpf: %w", exception.ErrInvalidCredentials)
	}
	token, err := r.token.FindByToken(ctx, input.Token)
	if err != nil {
		return fmt.Errorf("use case renew password, failed to find token: %w", err)
	}
	if token.IsExpired() {
		return fmt.Errorf("use case renew password, failed to renew password: %w", exception.ErrTokenExpired)
	}
	password := vo.Password(input.Password)
	if !password.IsValid() {
		return fmt.Errorf("use case renew password, invalid password: %w", exception.ErrInvalidPassword)
	}
	if err := r.user.UpdatePassword(ctx, token.IDUser, password.Value(), cpf.OnlyDigits()); err != nil {
		return fmt.Errorf("use case renew password, failed to update password: %w", err)
	}
	if err := r.token.MarkAsUsed(ctx, input.Token); err != nil {
		return fmt.Errorf("use case renew password, failed to mark token as used: %w", err)
	}
	return nil
}
