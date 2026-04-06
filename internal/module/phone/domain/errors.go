// Package domain
package domain

import "errors"

var (
	ErrInvalidID          = errors.New("invalid phone ID")
	ErrPhoneDeleted       = errors.New("phone is deleted")
	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrInvalidAreaCode    = errors.New("invalid area code")
	ErrInvalidKind        = errors.New("invalid kind")
	ErrInvalidNumber      = errors.New("invalid number")
	ErrInvalidDepartment  = errors.New("invalid department")
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
