// Package usecase
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()
	authentication := spy.NewAuthenticationSpy()
	logger := &spy.SpyLogger{}
	token := &spy.SpyToken{}
	smtp := &spy.SpySMTP{}
	sut := NewRecoverPassword(authentication, logger, token, smtp)
	input := dto.RecoverPasswordInput{Email: "email"}

	t.Run("should return error if email is invalid", func(t *testing.T) {
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Warn", ctx, "invalid email format provided", "email", input.Email).Once()
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, exception.ErrInvalidEmail)
	})

	t.Run("should return an error if it fails to find the id user", func(t *testing.T) {
		wantErr := errors.New("failed to find id user")
		authentication.FindError = wantErr
		input.Email = "email@example.com.br"
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Error", ctx, "account not found or is inactive", "error", authentication.FindError, "email", input.Email).Once()
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an errors if generate token its failed", func(t *testing.T) {
		wantErr := adapter.ErrBuildToken
		token.Err = wantErr
		authentication.FindError = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to generate reset token", "error", token.Err, "email", input.Email).Once()
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an error if it fails to create recover password", func(t *testing.T) {
		wantErr := errors.New("failed to create recover password")
		authentication.ValidateError = wantErr
		token.Err = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Info", ctx, "generated reset token", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to create reset token record", "error", authentication.ValidateError, "email", input.Email).Once()
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should return an error if it fails to send email", func(t *testing.T) {
		wantErr := errors.New("failed to send email")
		smtp.Err = wantErr
		authentication.ValidateError = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Info", ctx, "generated reset token", "email", input.Email).Once()
		logger.On("Info", ctx, "created reset token record", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to send email", "error", smtp.Err, "email", input.Email).Once()
		err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("should successful to recover password", func(t *testing.T) {
		smtp.Err = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Info", ctx, "generated reset token", "email", input.Email).Once()
		logger.On("Info", ctx, "created reset token record", "email", input.Email).Once()
		logger.On("Info", ctx, "recover password successful", "email", input.Email)
		err := sut.Execute(ctx, input)
		assert.NoError(t, err)
	})
}
