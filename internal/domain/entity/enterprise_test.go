// Package entity
package entity

import (
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/stretchr/testify/assert"
)

func TestEnterprise(t *testing.T) {
	t.Run("should return an enterprise", func(t *testing.T) {
		enterprise, err := NewEnterprise("Actajus", "Actajus", "10.123.456/0001-00")
		assert.NoError(t, err)
		assert.Equal(t, "Actajus", enterprise.Name.Value())
		assert.Equal(t, "Actajus", enterprise.TradeName.Value())
		assert.Equal(t, "10.123.456/0001-00", enterprise.CNPJ.Value())
		assert.True(t, enterprise.CreatedAt.Equal(enterprise.UpdatedAt))
		assert.True(t, time.Now().After(enterprise.CreatedAt))
		assert.True(t, time.Now().After(enterprise.UpdatedAt))
		assert.Nil(t, enterprise.DeletedAt)
	})

	t.Run("should return false if the name is invalid", func(t *testing.T) {
		_, err := NewEnterprise("11", "Actajus", "10.123.456/0001-00")
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidName, err)
	})

	t.Run("should return false if the trade name is invalid", func(t *testing.T) {
		_, err := NewEnterprise("Actajus", "11", "10.123.456/0001-00")
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidTradeName, err)
	})

	t.Run("should return false if the cnpj is invalid", func(t *testing.T) {
		_, err := NewEnterprise("Actajus", "Actajus", "10.123.456/0001")
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidCNPJ, err)
	})

	t.Run("should return false if enterprise is not deleted", func(t *testing.T) {
		enterprise, err := NewEnterprise("Actajus", "Actajus", "10.123.456/0001-00")
		assert.NoError(t, err)
		assert.False(t, enterprise.IsDeleted())
	})
}
