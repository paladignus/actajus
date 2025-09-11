// Package valueobject
package valueobject

import (
	"regexp"
	"strings"
	"unicode"
)

type CPF string

func (c CPF) Value() string {
	return string(c)
}

func (c CPF) IsValid() bool {
	re := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	if !re.MatchString(c.Value()) {
		return false
	}
	cpf := clear(c.Value())
	if allDigitsEqual(cpf) {
		return false
	}
	digit1 := calculateDigit(cpf[:9], 10)
	digit2 := calculateDigit(cpf[:10], 11)
	if digit1 != int(cpf[9]-'0') || digit2 != int(cpf[10]-'0') {
		return false
	}
	return true
}

func clear(cpf string) string {
	var builder strings.Builder
	for _, char := range cpf {
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

func calculateDigit(cpf string, factor int) int {
	sum := 0
	for _, char := range cpf {
		digit := int(char - '0')
		sum += digit * factor
		factor--
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}
