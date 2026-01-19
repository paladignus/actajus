// Package valueobject
package valueobject

import (
	"regexp"
)

type CPF string

func (c CPF) Value() string {
	return string(c)
}

func (c CPF) OnlyDigits() string {
	return clear(string(c))
}

func (c CPF) IsValid() bool {
	re := regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
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
