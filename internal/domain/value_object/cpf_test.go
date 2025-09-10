package valueobject

import (
	"regexp"
	"testing"
)

type CPF string

func (c CPF) Value() string {
	return string(c)
}

func (c CPF) IsValid() bool {
	re := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	return re.MatchString(c.Value())
}

func TestCPF(t *testing.T) {
	t.Run("should return false if cpf is invalid", func(t *testing.T) {
		sut := CPF("12345678900")
		if sut.IsValid() {
			t.Errorf("Expected CPF to be invalid, but got valid")
		}
	})
}

func (c CPF) formatted() string {
	re := regexp.MustCompile(`\.|-`)
	return re.ReplaceAllString(c.Value(), "")
}

func (c CPF) allDigitsEqual() bool {
	for i := 1; i < len(c.formatted()); i++ {
		if c.formatted()[i] != c.formatted()[0] {
			return false
		}
	}
	return true
}
