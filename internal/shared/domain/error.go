// Package domain
package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBlocked        = errors.New("user blocked")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionRevoked     = errors.New("session revoked")
	ErrSessionExpired     = errors.New("session expired")
	ErrRefreshMismatch    = errors.New("refresh token mismatch")
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors struct {
	errors []*FieldError
}

func NewValidationErrors(errors []*FieldError) *ValidationErrors {
	return &ValidationErrors{errors: errors}
}

func (e *ValidationErrors) Error() string {
	if len(e.errors) == 0 {
		return "validation errors"
	}
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}
	var messages []string
	for _, err := range e.errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf(
		"validation errors: %s",
		strings.Join(messages, "; "),
	)
}

func (e *ValidationErrors) Errors() []*FieldError {
	return e.errors
}
