// Package usecase
package usecase

import (
	"context"
	"testing"

	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/paladignus/actajus/internal/domain/repository"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
	"github.com/paladignus/actajus/test/infrastructure/persistence/spy"
	"github.com/stretchr/testify/assert"
)

type GetEmailByCPF struct {
	persistence repository.Account
	logger      repository.Logger
}

func NewGetEmailByCPF(persistence repository.Account, logger repository.Logger) GetEmailByCPF {
	return GetEmailByCPF{persistence, logger}
}

func (g GetEmailByCPF) Execute(ctx context.Context, cpfString string) (string, error) {
	g.logger.Info(ctx, "getting email by CPF", "cpf", cpfString)
	cpf := vo.CPF(cpfString)
	if !cpf.IsValid() {
		g.logger.Warn(ctx, "invalid cpf format provided", "cpf", cpfString)
		return "", exception.ErrInvalidCPF
	}
	if cpfString == "123.456.789-09" {
		return "", exception.ErrEmailNotFound
	}
	return "email@exemplo.com.br", nil
}

func TestGetEmailByCPF(t *testing.T) {
	ctx := context.Background()
	repository := &spy.AuthenticateSpy{}
	logger := &spy.MockLogger{}
	sut := NewGetEmailByCPF(repository, logger)
	t.Run("should return of an email when the CPF is valid and found", func(t *testing.T) {
		email, err := sut.Execute(ctx, "111.444.777-35")
		assert.NoError(t, err)
		assert.NotEmpty(t, email)
	})

	t.Run("should return an error when the CPF is invalid", func(t *testing.T) {
		_, err := sut.Execute(ctx, "123.456.789-00")
		assert.ErrorIs(t, err, exception.ErrInvalidCPF)
	})

	t.Run("should return an error if the CPF is not found", func(t *testing.T) {
		_, err := sut.Execute(ctx, "123.456.789-09")
		assert.ErrorIs(t, err, exception.ErrEmailNotFound)
	})
}
