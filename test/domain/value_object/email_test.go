package valueobject

import (
	"testing"

	valueobject "github.com/paladignus/actajus/internal/domain/value_object"
)

func TestEmail(t *testing.T) {
	t.Run("empty email should be invalid", func(t *testing.T) {
		sut := valueobject.Email("")
		if sut.IsValid() {
			t.Errorf("expected empty email to be invalid")
		}
	})
	t.Run("invalid email should be invalid", func(t *testing.T) {
		sut := valueobject.Email("invalid-email")
		if sut.IsValid() {
			t.Errorf("expected invalid email to be invalid")
		}
	})
	t.Run("equals method should work correctly", func(t *testing.T) {
		sut := valueobject.Email("marcelo@marcelo.eti.br")
		if !sut.Equals(valueobject.Email("marcelo@marcelo.eti.br")) {
			t.Errorf("expected emails to be equal")
		}
	})
}
