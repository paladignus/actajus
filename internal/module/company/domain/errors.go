// Package domain
package domain

import "errors"

var (
	ErrInvalidID           = errors.New("invalid company ID")
	ErrInvalidName         = errors.New("invalid company name")
	ErrInvalidTradeName    = errors.New("invalid trade name")
	ErrInvalidCNPJ         = errors.New("invalid CNPJ")
	ErrInvalidRegisteredBy = errors.New("invalid registered by")
	ErrCompanyDeleted      = errors.New("company is deleted")
	ErrAlreadyDeleted      = errors.New("company already deleted")
	ErrIDAlreadySet        = errors.New("company ID already set")
	ErrCompanyNotFound     = errors.New("company not found")
	ErrCNPJAlreadyExists   = errors.New("CNPJ already exists")
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
