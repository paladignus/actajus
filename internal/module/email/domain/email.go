// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Email struct {
	id        uint
	address   vo.Email
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type EmailBuilder struct {
	email  *Email
	errors []error
}

func NewEmailBuilder() *EmailBuilder {
	now := time.Now()
	return &EmailBuilder{
		email: &Email{
			createdAt: now,
			updatedAt: now,
		},
		errors: []error{},
	}
}

func (e *Email) UpdateBuilder() *EmailBuilder {
	if e.id == 0 {
		return &EmailBuilder{
			errors: []error{ErrInvalidID},
		}
	}
	if e.IsDeleted() {
		return &EmailBuilder{
			errors: []error{ErrEmailDeleted},
		}
	}
	return &EmailBuilder{
		email:  e,
		errors: []error{},
	}
}

func (e *EmailBuilder) WithID(id uint) *EmailBuilder {
	e.email.id = id
	if id == 0 {
		e.errors = append(e.errors, ErrInvalidID)
	}
	return e
}

func (e *EmailBuilder) WithAddress(address string) *EmailBuilder {
	e.email.address = vo.Email(address)
	if !e.email.address.IsValid() {
		e.errors = append(e.errors, ErrInvalidEmailAddress)
	}
	return e
}

func (e *EmailBuilder) WithCreatedAt(createdAt time.Time) *EmailBuilder {
	e.email.createdAt = createdAt
	return e
}

func (e *EmailBuilder) WithUpdatedAt(updatedAt time.Time) *EmailBuilder {
	e.email.updatedAt = updatedAt
	return e
}

func (e *EmailBuilder) WithDeletedAt(deletedAt *time.Time) *EmailBuilder {
	e.email.deletedAt = deletedAt
	return e
}

func (e *EmailBuilder) Build() (*Email, error) {
	if len(e.errors) > 0 {
		return nil, NewValidationErrors(e.errors)
	}
	if err := e.email.validate(); err != nil {
		return nil, err
	}
	return e.email, nil
}

func (e *EmailBuilder) Apply() error {
	if len(e.errors) > 0 {
		return NewValidationErrors(e.errors)
	}
	if err := e.email.validate(); err != nil {
		return err
	}
	e.email.updatedAt = time.Now()
	return nil
}

func (e *Email) ID() uint              { return e.id }
func (e *Email) Address() vo.Email     { return e.address }
func (e *Email) CreatedAt() time.Time  { return e.createdAt }
func (e *Email) UpdatedAt() time.Time  { return e.updatedAt }
func (e *Email) DeletedAt() *time.Time { return e.deletedAt }

func (e *Email) Delete() error {
	if e.id == 0 {
		return ErrInvalidID
	}
	if e.IsDeleted() {
		return ErrEmailDeleted
	}
	now := time.Now()
	e.deletedAt = &now
	e.updatedAt = now
	return nil
}

func (e *Email) IsDeleted() bool {
	return e.deletedAt != nil
}

func (e *Email) SetID(id uint) error {
	if e.id != 0 {
		return ErrInvalidID
	}
	e.id = id
	return nil
}

func (e *Email) validate() error {
	if !e.address.IsValid() {
		return ErrInvalidEmailAddress
	}
	return nil
}
