// Package usecase
package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestGetEmailByCPF(t *testing.T) {
	ctx := context.Background()
	authentication := &spy.Authentication{}
	logger := &spy.Logger{}
	input := dto.GetEmailByCPFInput{CPF: "111.444.777-35"}
	sut := NewGetEmailByCPF(authentication, logger)
	t.Run("should return of an email when the CPF is valid and found", func(t *testing.T) {
		authentication.FindResult.FindEmail.Email = "email@example.com.br"
		logger.On("Info", ctx, "getting email by CPF", "cpf", input.CPF).Once()
		logger.On("Info", ctx, "email found successfully", "cpf", input.CPF, "email", authentication.FindResult.FindEmail.Email).Once()
		email, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, email)
	})

	t.Run("should return an error when the CPF is invalid", func(t *testing.T) {
		input.CPF = "123.456.789-00"
		logger.On("Info", ctx, "getting email by CPF", "cpf", input.CPF).Once()
		logger.On("Warn", ctx, "invalid cpf format provided", "cpf", input.CPF).Once()
		_, err := sut.Execute(ctx, input)
		assert.ErrorIs(t, err, exception.ErrInvalidCPF)
	})

	t.Run("should return an error if the CPF is not found", func(t *testing.T) {
		input.CPF = "123.456.789-09"
		wantErr := exception.ErrEmailNotFound
		authentication.FindError = wantErr
		logger.On("Info", ctx, "getting email by CPF", "cpf", input.CPF).Once()
		logger.On("Warn", ctx, "email not found for provided cpf", "cpf", input.CPF, "error", wantErr).Once()
		_, err := sut.Execute(ctx, input)
		assert.ErrorIs(t, err, wantErr)
	})
}
