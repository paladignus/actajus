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
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

type RecoverPassword struct {
	persistence repository.Authentication
	logger      repository.Logger
	gateway     gateway.Token
}

func NewRecoverPassword(
	persistence repository.Authentication,
	logger repository.Logger,
	gateway gateway.Token,
) RecoverPassword {
	return RecoverPassword{persistence, logger, gateway}
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
		r.logger.Error(ctx, "failed to create recover password", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}
	r.logger.Info(ctx, "account is ative to email", "email", req.Email)

	token, err := r.gateway.GenerateResetToken(IDUser)
	if err != nil {
		r.logger.Error(ctx, "failed to create recover password", "error", err, "email", req.Email)
		return dto.RecoverPasswordOutput{}, err
	}

	r.repository.CreateRecoverPassword(ctx, IDUser, token)

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
	sut := NewRecoverPassword(authentication, logger, token)
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
		logger.On("Error", ctx, "failed to create recover password", "error", authentication.FindError, "email", input.Email).Once()
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, wantErr)
		assert.Empty(t, token.ResetToken)
	})

	t.Run("should successful to recover password", func(t *testing.T) {
		authentication.FindError = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "account is ative to email", "email", input.Email).Once()
		logger.On("Info", ctx, "recover password successful", "email", input.Email)
		token, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, token.RecoverToken)
	})
}
