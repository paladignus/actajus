// Package domain
package domain

import "errors"

var (
	ErrInvalidID          = errors.New("invalid social media ID")
	ErrSocialMediaDeleted = errors.New("social media is deleted")
	ErrInvalidIDCompany   = errors.New("invalid id company")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidURL         = errors.New("invalid url")
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
