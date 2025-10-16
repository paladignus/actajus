// Package valueobject
package valueobject

import (
	"testing"
)

func TestEmail(t *testing.T) {
	t.Run("empty email should be invalid", func(t *testing.T) {
		sut := Email("")
		if sut.IsValid() {
			t.Errorf("expected empty email to be invalid")
		}
	})
	t.Run("invalid email should be invalid", func(t *testing.T) {
		sut := Email("invalid-email")
		if sut.IsValid() {
			t.Errorf("expected invalid email to be invalid")
		}
	})
	t.Run("equals method should work correctly", func(t *testing.T) {
		sut := Email("marcelo@marcelo.eti.br")
		if !sut.Equals(Email("marcelo@marcelo.eti.br")) {
			t.Errorf("expected emails to be equal")
		}
	})
}
