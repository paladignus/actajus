// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Email struct {
	id        int64
	address   vo.Email
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type EmailBuilder struct {
	email *Email
}

func NewEmailBuilder() *EmailBuilder {
	now := time.Now()
	return &EmailBuilder{
		email: &Email{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (e *EmailBuilder) WithID(id int64) *EmailBuilder {
	e.email.id = id
	return e
}

func (e *EmailBuilder) WithAddress(address string) *EmailBuilder {
	e.email.address = vo.Email(address)
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
	if err := e.email.validate(); err != nil {
		return nil, err
	}
	return e.email, nil
}

func (e *EmailBuilder) Apply() error {
	if err := e.email.validate(); err != nil {
		return err
	}
	e.email.updatedAt = time.Now()
	return nil
}

func (e *Email) ID() int64             { return e.id }
func (e *Email) Address() vo.Email     { return e.address }
func (e *Email) CreatedAt() time.Time  { return e.createdAt }
func (e *Email) UpdatedAt() time.Time  { return e.updatedAt }
func (e *Email) DeletedAt() *time.Time { return e.deletedAt }

func (e *Email) Delete() error {
	if e.id == 0 {
		return domain.NewFieldError("id", "email ID is invalid")
	}
	if e.IsDeleted() {
		return domain.NewFieldError("deleted_at", "email is already deleted")
	}
	now := time.Now()
	e.deletedAt = &now
	e.updatedAt = now
	return nil
}

func (e *Email) IsDeleted() bool {
	return e.deletedAt != nil
}

func (e *Email) SetID(id int64) error {
	if e.id != 0 {
		return domain.NewFieldError("id", "email ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "email ID is invalid")
	}
	e.id = id
	return nil
}

func (e *Email) validate() error {
	if !e.address.IsValid() {
		return domain.NewFieldError("address", "email address is invalid")
	}
	return nil
}
