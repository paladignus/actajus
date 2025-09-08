package valueobject

import (
	"regexp"
	"testing"
)

type Password string

func (p Password) Value() string {
	return string(p)
}

func (p Password) IsEmpty() bool {
	return len(p) == 0
}

func (p Password) IsValid() bool {
	number := regexp.MustCompile(`[0-9]`)
	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	special := regexp.MustCompile(`[\W_]`)
	return !p.IsEmpty() &&
		len(p.Value()) > 7 &&
		number.MatchString(p.Value()) &&
		lower.MatchString(p.Value()) &&
		upper.MatchString(p.Value()) &&
		special.MatchString(p.Value())
}

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
