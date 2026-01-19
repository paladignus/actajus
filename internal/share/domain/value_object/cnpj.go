// Package valueobject
package valueobject

import "regexp"

type CNPJ string

func (c CNPJ) Value() string {
	return string(c)
}

func (c CNPJ) OnlyDigits() string {
	return clear(string(c))
}

func (c CNPJ) IsValid() bool {
	re := regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}$`)
	if !re.MatchString(c.Value()) {
		return false
	}
	cnpj := clear(c.Value())
	if allDigitsEqual(cnpj) {
		return false
	}
	digit1 := calculateDigit(cnpj[:12], 5)
	digit2 := calculateDigit(cnpj[:13], 6)
	if digit1 != int(cnpj[12]-'0') || digit2 != int(cnpj[13]-'0') {
		return false
	}
	return true
}
