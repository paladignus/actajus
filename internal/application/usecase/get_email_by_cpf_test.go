// Package usecase
package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/infrastructure/persistence/spy"
	"github.com/stretchr/testify/assert"
)

func TestGetEmailByCPF(t *testing.T) {
	ctx := context.Background()
	account := &spy.AuthenticateSpy{}
	logger := &spy.MockLogger{}
	req := dto.GetEmailByCPFInput{CPF: "111.444.777-35"}
	sut := NewGetEmailByCPF(account, logger)
	t.Run("should return of an email when the CPF is valid and found", func(t *testing.T) {
		account.FindResult.Email = "email@example.com.br"
		email, err := sut.Execute(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, email)
	})

	t.Run("should return an error when the CPF is invalid", func(t *testing.T) {
		req.CPF = "123.456.789-00"
		_, err := sut.Execute(ctx, req)
		assert.ErrorIs(t, err, exception.ErrInvalidCPF)
	})

	t.Run("should return an error if the CPF is not found", func(t *testing.T) {
		req.CPF = "123.456.789-09"
		wantErr := exception.ErrEmailNotFound
		account.FindError = wantErr
		_, err := sut.Execute(ctx, req)
		assert.ErrorIs(t, err, wantErr)
	})
}
