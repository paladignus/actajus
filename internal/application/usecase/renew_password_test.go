// Package usecase
package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestRenewPassword(t *testing.T) {
	ctx := context.Background()
	user := &spy.User{}
	token := &spy.TokenSpy{}
	sut := NewRenewPassword(user, token)
	input := dto.RenewPasswordInput{
		CPF: "invalid-cpf",
	}

	t.Run("should return an erro if cpf is invalid", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, invalid cpf: invalid credentials")
		assert.ErrorIs(t, err, exception.ErrInvalidCredentials)
	})

	input.CPF = "987.654.321-00"
	t.Run("should return an error if token not found", func(t *testing.T) {
		expectErr := errors.New("token not found")
		token.FindErr = expectErr
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, failed to find token: token not found")
		assert.ErrorIs(t, err, expectErr)
	})

	expiredToken := entity.PasswordResetToken{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	token.FindErr = nil
	t.Run("should return an error if token is expired", func(t *testing.T) {
		token.FindResult = expiredToken
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, failed to renew password: token expired")
		assert.ErrorIs(t, err, exception.ErrTokenExpired)
	})

	validToken := entity.PasswordResetToken{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	token.FindResult = validToken
	t.Run("should return an error if password is invalid", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, invalid password: invalid password")
		assert.ErrorIs(t, err, exception.ErrInvalidPassword)
	})

	input.Password = "Valid@123"
	user.ValidateError = errors.New("update password error")
	t.Run("should return an error if update password fails", func(t *testing.T) {
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, failed to update password: update password error")
		assert.ErrorIs(t, err, user.ValidateError)
	})

	t.Run("should return an error if mark token as used fails", func(t *testing.T) {
		user.ValidateError = nil
		token.MaskErr = errors.New("mark as used error")
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "use case renew password, failed to mark token as used: mark as used error")
		assert.ErrorIs(t, err, token.MaskErr)
	})

	t.Run("should return no error", func(t *testing.T) {
		token.MaskErr = nil
		err := sut.Execute(ctx, input)
		assert.NoError(t, err)
	})

	// t.Run("should return error token is invalid", func(t *testing.T) {
	// 	sut := NewRenewPassword(nil)
	// 	err := sut.Execute(ctx, "token")
	// 	assert.Error(t, err)
	// 	assert.ErrorContains(t, err, "invalid token")
	// })
}
