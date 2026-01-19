// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	t.Run("empty email should be invalid", func(t *testing.T) {
		sut := Email("")
		assert.False(t, sut.IsValid(), "expected empty email to be invalid")
	})
	t.Run("invalid email should be invalid", func(t *testing.T) {
		sut := Email("invalid-email")
		assert.False(t, sut.IsValid(), "expected invalid email to be invalid")
	})
	t.Run("equals method should work correctly", func(t *testing.T) {
		sut := Email("marcelo@marcelo.eti.br")
		assert.True(t, sut.Equals(Email("marcelo@marcelo.eti.br")), "expected emails to be equal")
	})
}
