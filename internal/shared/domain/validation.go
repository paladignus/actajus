// Package domain
package domain

type ViolationCode string

const (
	CodeRequired   ViolationCode = "required"
	CodeInvalid    ViolationCode = "invalid"
	CodeMin        ViolationCode = "min"
	CodeMax        ViolationCode = "max"
	CodeLen        ViolationCode = "len"
	CodeNumeric    ViolationCode = "numeric"
	CodeEmail      ViolationCode = "email"
	CodeOneOf      ViolationCode = "oneof"
	CodeRequiredIf ViolationCode = "required_if"
)

type Violation struct {
	Path string            `json:"path"`
	Code ViolationCode     `json:"code"`
	Meta map[string]string `json:"meta,omitzero"`
}

type ValidationError struct {
	Violations []Violation `json:"violations"`
}

func NewValidationError(violations []Violation) error {
	if len(violations) == 0 {
		return nil
	}
	return ValidationError{Violations: violations}
}

func (e ValidationError) Error() string { return "validation_error" }
