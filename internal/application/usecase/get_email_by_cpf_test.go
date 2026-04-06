// Package usecase
package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/test/spy"
	"github.com/stretchr/testify/assert"
)

func TestGetEmailByCPF(t *testing.T) {
	ctx := context.Background()
	user := &spy.User{}
	input := command.GetEmailByCPFCommand{CPF: "111.444.777-35"}
	sut := NewGetEmailByCPF(user)
	t.Run("should return of an email when the CPF is valid and found", func(t *testing.T) {
		user.FindResult.FindEmail.Email = "email@example.com.br"
		email, err := sut.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, email)
	})

	t.Run("should return an error when the CPF is invalid", func(t *testing.T) {
		input.CPF = "123.456.789-00"
		_, err := sut.Execute(ctx, input)
		assert.Error(t, err)
		assert.ErrorIs(t, err, exception.ErrEmailNotFound)
	})

	t.Run("should return an error if the CPF is not found", func(t *testing.T) {
		input.CPF = "123.456.789-09"
		wantErr := exception.ErrEmailNotFound
		user.FindError = wantErr
		_, err := sut.Execute(ctx, input)
		assert.ErrorIs(t, err, wantErr)
	})
}
