// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	sut := Text("M4rc3l0@#$%¨&*()")
	t.Run("should return error if text contains numbers and special characters", func(t *testing.T) {
		assert.False(t, sut.IsValid(), "Expected error for text with numbers and special characters, but got nil")
	})
	sut = Text("Marcelo Bento Pereira")
	t.Run("should return nil if text is valid", func(t *testing.T) {
		assert.True(t, sut.IsValid(), "Expected nil for valid text, but got error")
	})
}
