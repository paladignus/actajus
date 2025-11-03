// Package usecase
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/gateway"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

type RecoverPassword struct {
	persistence repository.Authentication
	logger      repository.Logger
	token       gateway.Token
	smtp        gateway.SMTP
}

func NewRecoverPassword(
	persistence repository.Authentication,
	logger repository.Logger,
	token gateway.Token,
	smtp gateway.SMTP,
) RecoverPassword {
	return RecoverPassword{persistence, logger, token, smtp}
}

func (r RecoverPassword) Execute(ctx context.Context, req dto.RecoverPasswordInput) (dto.RecoverPasswordOutput, error) {
	r.logger.Info(ctx, "recover password", "email", req.Email)
	email := vo.Email(req.Email)
	if !email.IsValid() {
		r.logger.Warn(ctx, "invalid email format provided",
			"email", req.Email,
		)
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidEmail
	}

	IDUser, err := r.persistence.AccountIsActive(ctx, req.Email)
	if err != nil {
		r.logger.Error(ctx, "account not found or is inactive", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "account is ative to email", "email", req.Email)

	token, err := r.token.GenerateResetToken(IDUser)
	if err != nil {
		r.logger.Error(ctx, "failed to generate reset token", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	if err = r.persistence.CreateRecoverPassword(ctx, IDUser, token.ResetToken); err != nil {
		r.logger.Error(ctx, "failed to create reset token record", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	if err = r.smtp.SendEmail(ctx, email.Value(), token.ResetToken); err != nil {
		r.logger.Error(ctx, "failed to send email", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.logger.Info(ctx, "recover password successful", "email", req.Email)
	return dto.RecoverPasswordOutput{
		RecoverToken: "token",
	}, nil
}

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
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, exception.ErrInvalidEmail)
	})

	t.Run("should return an error if it fails to find the id user", func(t *testing.T) {
		wantErr := errors.New("failed to find id user")
		authentication.FindError = wantErr
		input.Email = "email@example.com.br"
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Error", ctx, "account not found or is inactive", "error", authentication.FindError, "email", input.Email).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, token.ResetToken)
	})

	t.Run("should return an errors if generate token its failed", func(t *testing.T) {
		wantErr := adapter.ErrBuildToken
		token.Err = wantErr
		authentication.FindError = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to generate reset token", "error", token.Err, "email", input.Email).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, token.ResetToken)
	})

	t.Run("should return an error if it fails to create recover password", func(t *testing.T) {
		wantErr := errors.New("failed to create recover password")
		authentication.ValidateError = wantErr
		token.Err = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to create reset token record", "error", authentication.ValidateError, "email", input.Email).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, token.ResetToken)
	})

	t.Run("should return an error if it fails to send email", func(t *testing.T) {
		wantErr := errors.New("failed to send email")
		smtp.Err = wantErr
		authentication.ValidateError = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Error", ctx, "failed to send email", "error", smtp.Err, "email", input.Email).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, token.ResetToken)
	})

	t.Run("should successful to recover password", func(t *testing.T) {
		smtp.Err = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Info", ctx, "recover password successful", "email", input.Email)
		token, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, token.RecoverToken)
	})
}
