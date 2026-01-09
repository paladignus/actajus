// Package dto
package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnterpriseInput(t *testing.T) {
	sut := EnterpriseInput{
		RegisteredBy: 1,
		Name:         "Teste",
		TradeName:    "Teste Trade",
		CNPJ:         "10.123.456/0001-00",
	}

	t.Run("should return true same values", func(t *testing.T) {
		assert.Equal(t, sut.RegisteredBy, 1)
		assert.Equal(t, sut.Name, "Teste")
		assert.Equal(t, sut.TradeName, "Teste Trade")
		assert.Equal(t, sut.CNPJ, "10.123.456/0001-00")
	})

	t.Run("should return empty values", func(t *testing.T) {
		sut := EnterpriseInput{}
		assert.Equal(t, sut.RegisteredBy, 0)
		assert.Equal(t, sut.Name, "")
		assert.Equal(t, sut.TradeName, "")
		assert.Equal(t, sut.CNPJ, "")
	})
}
