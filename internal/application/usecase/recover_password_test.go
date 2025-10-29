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
	if ok := vo.CPF(req.CPF).IsValid(); !ok {
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidCPF
	}
	if ok := vo.Email(req.Email).IsValid(); !ok {
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidEmail
	}
	if ok := vo.Password(req.Password).IsValid(); !ok {
		return dto.RecoverPasswordOutput{}, exception.ErrInvalidPassword
	}
	if r.ErrorToken != nil {
		return dto.RecoverPasswordOutput{}, r.ErrorToken
	}
	return dto.RecoverPasswordOutput{
		RecoverToken: "token",
	}, nil
}

func TestRecoverPassword(t *testing.T) {
	ctx := context.Background()
	logger := &spy.SpyLogger{}
	sut := NewRecoverPassword(logger)
	input := dto.RecoverPasswordInput{}

	t.Run("should return error if cpf is invalid", func(t *testing.T) {
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, exception.ErrInvalidCPF)
	})

	t.Run("should return error if email is invalid", func(t *testing.T) {
		input.CPF = "111.444.777-35"
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, exception.ErrInvalidEmail)
	})

	t.Run("should return error if password is invalid", func(t *testing.T) {
		input.Email = "email@example.com.br"
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, exception.ErrInvalidPassword)
	})

	ErrGenerateToken := errors.New("failed to generate token")

	t.Run("should return an error if it fails to generate the token", func(t *testing.T) {
		input.Password = "123456@Ma"
		sut.ErrorToken = ErrGenerateToken
		token, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.Empty(t, token.RecoverToken)
		assert.ErrorIs(t, err, ErrGenerateToken)
	})

	t.Run("should return recover token", func(t *testing.T) {
		sut.ErrorToken = nil
		token, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, token.RecoverToken)
	})
}
