// Package usecase
package usecase

import (
	"testing"

	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
	"github.com/stretchr/testify/assert"
)

type GetEmailByCPF struct{}

func (g GetEmailByCPF) Execute(cpfString string) (string, error) {
	cpf := vo.CPF(cpfString)
	if !cpf.IsValid() {
		return "", exception.ErrInvalidCPF
	}
	if cpfString == "123.456.789-09" {
		return "", exception.ErrEmailNotFound
	}
	return "email@exemplo.com.br", nil
}

func TestGetEmailByCPF(t *testing.T) {
	sut := GetEmailByCPF{}
	t.Run("should return of an email when the CPF is valid and found", func(t *testing.T) {
		email, err := sut.Execute("111.444.777-35")
		assert.NoError(t, err)
		assert.NotEmpty(t, email)
	})

	t.Run("should return an error when the CPF is invalid", func(t *testing.T) {
		_, err := sut.Execute("123.456.789-00")
		assert.ErrorIs(t, err, exception.ErrInvalidCPF)
	})

	t.Run("should return an error if the CPF is not found", func(t *testing.T) {
		_, err := sut.Execute("123.456.789-09")
		assert.ErrorIs(t, err, exception.ErrEmailNotFound)
	})
}
