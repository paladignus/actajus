// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCNPJ(t *testing.T) {
	t.Run("should return false if cpf is invalid", func(t *testing.T) {
		invalidCnpjs := []CNPJ{
			"11.111.111/1111-11",  // Inválido (todos dígitos iguais)
			"10.123.456/0001-09",  // Inválido
			"1012345600010",       // Inválido (menos de 14 dígitos)
			"10.123.456/0001-000", // Inválido (mais de 14 dígitos)
			"10.123.456/0001-XX",  // Inválido (contém letras)
		}
		for _, cnpj := range invalidCnpjs {
			assert.Falsef(t, cnpj.IsValid(), "Expected CCNPJ to be invalid, but got valid: %s", cnpj.Value())
		}
	})

	t.Run("should return true if cnpj is valid", func(t *testing.T) {
		validCnpjs := []CNPJ{
			"10.123.456/0001-00",
			"10123456000100",
			"25.461.378/0001-13",
			"25461378000113",
			"66.628.623/0001-11",
			"66628623000111",
		}
		for _, cnpj := range validCnpjs {
			assert.Truef(t, cnpj.IsValid(), "Expected CCNPJ to be valid, but got invalid: %s", cnpj.Value())
		}
	})

	t.Run("should return the only digits", func(t *testing.T) {
		validCnpjs := []CNPJ{
			"10.123.456/0001-00",
		}
		for _, cnpj := range validCnpjs {
			expect := cnpj.OnlyDigits()
			assert.Equal(t, "10123456000100", expect, "Expected only digits to be '10123456000100', but got %s", expect)
		}
	})
}
