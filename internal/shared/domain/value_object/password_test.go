// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	sut := Password("")
	t.Run("empty password should be invalid", func(t *testing.T) {
		assert.False(t, sut.IsValid(), "expected empty password to be invalid")
	})
	t.Run("non-empty password should be valid", func(t *testing.T) {
		sut = Password("Securepassword2!@#$%&*()_+/?;:.><,~^")
		assert.True(t, sut.IsValid(), "expected non-empty password to be valid")
	})
}
