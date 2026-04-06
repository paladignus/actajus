// Package domain
package domain

import "errors"

var (
	ErrInvalidID           = errors.New("invalid email ID")
	ErrEmailDeleted        = errors.New("email is deleted")
	ErrInvalidEmailAddress = errors.New("invalid email address")
)

type ValidationErrors struct {
	Errors []error
}

func NewValidationErrors(errors []error) *ValidationErrors {
	return &ValidationErrors{Errors: errors}
}

func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation errors"
	}
	return e.Errors[0].Error()
}
