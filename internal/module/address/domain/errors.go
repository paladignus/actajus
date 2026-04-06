// Package domain
package domain

import "errors"

var (
	ErrInvalidID           = errors.New("invalid address ID")
	ErrInvalidZIP          = errors.New("invalid ZIP code")
	ErrInvalidTitle        = errors.New("invalid title")
	ErrInvalidStreet       = errors.New("invalid street")
	ErrInvalidNumber       = errors.New("invalid number")
	ErrInvalidNeighborhood = errors.New("invalid neighborhood")
	ErrInvalidCity         = errors.New("invalid city")
	ErrInvalidState        = errors.New("invalid state")
	ErrInvalidCountry      = errors.New("invalid country")
	ErrAddressDeleted      = errors.New("address is deleted")
	ErrAlreadyDeleted      = errors.New("address already deleted")
	ErrIDAlreadySet        = errors.New("address ID already set")
	ErrAddressNotFound     = errors.New("address not found")
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
