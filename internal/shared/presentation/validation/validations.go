// Package validation
package validation

type Code string

const (
	CodeRequired   Code = "required"
	CodeMin        Code = "min"
	CodeMax        Code = "max"
	CodeLen        Code = "len"
	CodeNumeric    Code = "numeric"
	CodeOneOf      Code = "oneof"
	CodeEmail      Code = "email"
	CodeRequiredIf Code = "required_if"
)

type Violation struct {
	Path string
	Code Code
	Meta map[string]string
}
