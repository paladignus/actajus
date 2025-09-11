// Package valueobject
package valueobject

import (
	"regexp"
	"testing"
)

type Date string

func (d Date) Value() string {
	return string(d)
}

func (d Date) IsValid() bool {
	re := regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$`)
	return re.MatchString(d.Value())
}

func TestDate(t *testing.T) {
	t.Run("should return false if date is invalid", func(t *testing.T) {
		datesInvalids := []Date{
			"32/01/2020", // Dia inválido
			"00/01/2020", // Dia inválido
			"15/13/2020", // Mês inválido
			"15/00/2020", // Mês inválido
			"15/01/20",   // Ano inválido
			"15-01-2020", // Formato inválido
			"15012020",   // Formato inválido
		}
		for _, date := range datesInvalids {
			if date.IsValid() {
				t.Errorf("Expected Date to be invalid, but got valid")
			}
		}
	})
	t.Run("should return true if date is valid", func(t *testing.T) {
		datesValids := []Date{
			"01/01/2020",
			"15/06/1995",
			"31/12/2023",
		}
		for _, date := range datesValids {
			if !date.IsValid() {
				t.Errorf("Expected Date to be valid, but got invalid")
			}
		}
	})
}
