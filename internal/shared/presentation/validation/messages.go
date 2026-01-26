// Package validation
package validation

type Lang string

const (
	PT Lang = "pt"
	EN Lang = "en"
)

var messages = map[Lang]map[string]string{
	PT: {
		"required": "%s é obrigatório",
		"min":      "%s deve ser no mínimo %d",
		"len":      "%s deve conter exatamente %d caracteres",
		"numeric":  "%s deve conter apenas números",
	},
	EN: {
		"required": "%s is required",
		"min":      "%s must be at least %d",
		"len":      "%s must be exactly %d characters",
		"numeric":  "%s must contain only numbers",
	},
}
