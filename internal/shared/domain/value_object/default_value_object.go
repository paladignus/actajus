// Package valueobject
package valueobject

import (
	"strings"
	"unicode"
)

func clear(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func allDigitsEqual(cpf string) bool {
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			return false
		}
	}
	return true
}

func calculateDigit(input string, factor int) int {
	var isCNPJ bool
	if len(input) > 11 {
		isCNPJ = true
	}
	sum := 0
	for _, char := range input {
		digit := int(char - '0')
		sum += digit * factor
		factor--
		if isCNPJ && factor < 2 {
			factor = 9
		}
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}
