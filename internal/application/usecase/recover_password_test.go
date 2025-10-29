// Package usecase
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

type IRecoverPassword interface {
	Execute(ctx context.Context, req dto.RecoverPasswordInput) (dto.RecoverPasswordOutput, error)
}

type RecoverPassword struct {
	logger     repository.Logger
	ErrorToken error
}

func NewRecoverPassword(logger repository.Logger) RecoverPassword {
	return RecoverPassword{logger, nil}
}

func (r RecoverPassword) Execute(ctx context.Context, req dto.RecoverPasswordInput) (dto.RecoverPasswordOutput, error) {
	r.logger.Info(ctx, "recover password", "email", req.Email)
	if ok := vo.Email(req.Email).IsValid(); !ok {
		r.logger.Warn(ctx, "invalid email format provided", "email", req.Email)
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidEmail
	}
	if r.ErrorToken != nil {
		r.logger.Warn(ctx, "failed to generate token", "error", r.ErrorToken)
		return dto.RecoverPasswordOutput{}, r.ErrorToken
	}
	r.logger.Info(ctx, "recover password successful", "email", req.Email)
	return dto.RecoverPasswordOutput{
		RecoverToken: "token",
	}, nil
}

func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()
	logger := &spy.SpyLogger{}
	sut := NewRecoverPassword(logger)
	input := dto.RecoverPasswordInput{Email: "email"}

	t.Run("should return error if email is invalid", func(t *testing.T) {
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Warn", ctx, "invalid email format provided", "email", input.Email).Once()
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, exception.ErrInvalidEmail)
	})

	t.Run("should return an error if it fails to generate the token", func(t *testing.T) {
		input.Email = "email@example.com.br"
		ErrGenerateToken := errors.New("failed to generate token")
		sut.ErrorToken = ErrGenerateToken
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Warn", ctx, "failed to generate token", "error", sut.ErrorToken).Once()
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, ErrGenerateToken)
	})

	t.Run("should return recover token", func(t *testing.T) {
		sut.ErrorToken = nil
		logger.On("Info", ctx, "recover password", "email", input.Email).Once()
		logger.On("Info", ctx, "recover password successful", "email", input.Email).Once()
		token, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, token.RecoverToken)
	})
}
