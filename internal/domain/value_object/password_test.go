// Package valueobject
package valueobject

import (
	"testing"
)

func TestPassword(t *testing.T) {
	sut := Password("")
	t.Run("empty password should be invalid", func(t *testing.T) {
		if sut.IsValid() {
			t.Errorf("expected empty password to be invalid")
		}
	})
	t.Run("non-empty password should be valid", func(t *testing.T) {
		sut = Password("Securepassword2!@#$%&*()_+/?;:.><,~^")
		if !sut.IsValid() {
			t.Errorf("expected non-empty password to be valid")
		}
	})
}
