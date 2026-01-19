// Package valueobject
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValueObject(t *testing.T) {
	t.Run("should return only digits", func(t *testing.T) {
		sut := clear("12.345.678-90")
		assert.Equal(t, "1234567890", sut, "Expected only digits, but got %s", sut)
	})

	t.Run("should return empty string", func(t *testing.T) {
		sut := clear("")
		assert.Equal(t, "", sut, "Expected empty string, but got %s", sut)
	})

	t.Run("should return true if all characters is equals", func(t *testing.T) {
		sut := allDigitsEqual("11111111111")
		assert.True(t, sut, "Expected true, but got false")
	})

	t.Run("should return false if all characters is not equals", func(t *testing.T) {
		sut := allDigitsEqual("11111111112")
		assert.False(t, sut, "Expected false, but got true")
	})

	t.Run("should calculate the second-to-last and last digits", func(t *testing.T) {
		cpf := "52998224725"
		sut := calculateDigit(cpf[:9], 10)
		assert.Equal(t, int(cpf[9]-'0'), sut, "Expected 0, but got %d", sut)
		sut = calculateDigit(cpf[:10], 11)
		assert.Equal(t, int(cpf[10]-'0'), sut, "Expected 0, but got %d", sut)
		cnpj := "10123456000100"
		sut = calculateDigit(cnpj[:12], 5)
		assert.Equal(t, int(cnpj[12]-'0'), sut, "Expected 0, but got %d", sut)
		sut = calculateDigit(cnpj[:13], 6)
		assert.Equal(t, int(cnpj[13]-'0'), sut, "Expected 0, but got %d", sut)
	})
}
