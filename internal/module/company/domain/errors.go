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
	Errs []error
}

func NewValidationErrors(errs []error) *ValidationErrors {
	return &ValidationErrors{Errs: errs}
}

func (e *ValidationErrors) Error() string {
	if len(e.Errs) == 0 {
		return "validation errors"
	}
	return e.Errs[0].Error()
}

func (e *ValidationErrors) Errors() []error {
	return e.Errs
}
