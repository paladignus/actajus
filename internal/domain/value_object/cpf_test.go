// Package valueobject
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
			"72775351115",
		}
		for _, cpf := range validCpfs {
			if !cpf.IsValid() {
				t.Errorf("Expected CPF to be valid, but got invalid")
			}
		}
	})
	t.Run("", func(t *testing.T) {
		validCpfs := []CPF{
			"529.982.247-25",
		}
		for _, cpf := range validCpfs {
			expect := cpf.OnlyDigits()
			if expect != "52998224725" {
				t.Errorf("Expected only digits to be '52998224725', but got %s", expect)
			}
		}
	})
}
