// Package entity
package entity

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/stretchr/testify/assert"
)

func TestEnterprise(t *testing.T) {
	input := dto.Enterprise{
		RegisteredBy: 1, Name: "Actajus", TradeName: "Actajus", CNPJ: "10.123.456/0001-00",
	}
	t.Run("should return an enterprise", func(t *testing.T) {
		enterprise, err := NewEnterprise(input)
		assert.NoError(t, err)
		assert.Equal(t, "Actajus", enterprise.Name.Value())
		assert.Equal(t, "Actajus", enterprise.TradeName.Value())
		assert.Equal(t, "10.123.456/0001-00", enterprise.CNPJ.Value())
		assert.True(t, enterprise.CreatedAt.Equal(enterprise.UpdatedAt))
		assert.True(t, time.Now().After(enterprise.CreatedAt))
		assert.True(t, time.Now().After(enterprise.UpdatedAt))
		assert.Nil(t, enterprise.DeletedAt)
	})

	t.Run("should return false if the registered by is invalid", func(t *testing.T) {
		input.RegisteredBy = 0
		_, err := NewEnterprise(input)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidRegisteredBy, err)
	})

	t.Run("should return false if the name is invalid", func(t *testing.T) {
		input.RegisteredBy = 1
		input.Name = "11"
		_, err := NewEnterprise(input)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidName, err)
	})

	t.Run("should return false if the trade name is invalid", func(t *testing.T) {
		input.Name = "Actajus"
		input.TradeName = "11"
		_, err := NewEnterprise(input)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidTradeName, err)
	})

	t.Run("should return false if the cnpj is invalid", func(t *testing.T) {
		input.TradeName = "Trade Actajus"
		input.CNPJ = "10.123.456/0001"
		_, err := NewEnterprise(input)
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidCNPJ, err)
	})

	t.Run("should return false if enterprise is not deleted", func(t *testing.T) {
		input.CNPJ = "10.123.456/0001-00"
		enterprise, err := NewEnterprise(input)
		assert.NoError(t, err)
		assert.False(t, enterprise.IsDeleted())
	})
}
