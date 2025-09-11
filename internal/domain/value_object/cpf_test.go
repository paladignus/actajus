package valueobject

import (
	"testing"
)

func TestCPF(t *testing.T) {
	t.Run("should return false if cpf is invalid", func(t *testing.T) {
		invalidCpfs := []CPF{
			"111.111.111-11", // Inválido (todos dígitos iguais)
			"123.456.789-00", // Inválido
			"529982247",      // Inválido (menos de 11 dígitos)
			"529.982.247-XX", // Inválido (contém letras)
		}
		for _, cpf := range invalidCpfs {
			if cpf.IsValid() {
				t.Errorf("Expected CPF to be invalid, but got valid")
			}
		}
	})
	t.Run("should return true if cpf is valid", func(t *testing.T) {
		validCpfs := []CPF{
			"529.982.247-25",
			"727.753.511-15",
			"932.272.131-68",
		}
		for _, cpf := range validCpfs {
			if !cpf.IsValid() {
				t.Errorf("Expected CPF to be valid, but got invalid")
			}
		}
	})
}
