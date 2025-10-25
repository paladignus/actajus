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
	return "email@examplo.com.br", nil
}

func TestGetEmailByCPF(t *testing.T) {
	sut := GetEmailByCPF{}
	t.Run("should return of an email when the CPF is valid", func(t *testing.T) {
		_, err := sut.Execute("111.444.777-35")
		assert.NoError(t, err)
	})
}
