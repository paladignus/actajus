// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Phone struct {
	id         uint
	number     vo.PhoneNumber
	kind       vo.Text
	department vo.Text
	createdAt  time.Time
	updatedAt  time.Time
	deletedAt  *time.Time
}

type PhoneBuilder struct {
	phone  *Phone
	errors []error
}

func NewPhoneBuilder() *PhoneBuilder {
	now := time.Now()
	return &PhoneBuilder{
		phone: &Phone{
			createdAt: now,
			updatedAt: now,
		},
		errors: []error{},
	}
}

func (p *Phone) UpdateBuilder() *PhoneBuilder {
	if p.id == 0 {
		return &PhoneBuilder{
			errors: []error{ErrInvalidID},
		}
	}
	if p.IsDeleted() {
		return &PhoneBuilder{
			errors: []error{ErrPhoneDeleted},
		}
	}
	return &PhoneBuilder{
		phone:  p,
		errors: []error{},
	}
}

func (p *PhoneBuilder) WithID(id uint) *PhoneBuilder {
	p.phone.id = id
	if id == 0 {
		p.errors = append(p.errors, ErrInvalidID)
	}
	return p
}

func (p *PhoneBuilder) WithNumber(number string) *PhoneBuilder {
	p.phone.number = vo.PhoneNumber(number)
	if !p.phone.number.IsValid() {
		p.errors = append(p.errors, ErrInvalidNumber)
	}
	return p
}

func (p *PhoneBuilder) WithKind(kind string) *PhoneBuilder {
	p.phone.kind = vo.Text(kind)
	if !p.phone.kind.IsValid() {
		p.errors = append(p.errors, ErrInvalidKind)
	}
	return p
}

func (p *PhoneBuilder) WithDepartment(department string) *PhoneBuilder {
	p.phone.department = vo.Text(department)
	if !p.phone.department.IsValid() {
		p.errors = append(p.errors, ErrInvalidDepartment)
	}
	return p
}

func (p *PhoneBuilder) WithCreatedAt(createdAt time.Time) *PhoneBuilder {
	p.phone.createdAt = createdAt
	return p
}

func (p *PhoneBuilder) WithUpdatedAt(updatedAt time.Time) *PhoneBuilder {
	p.phone.updatedAt = updatedAt
	return p
}

func (p *PhoneBuilder) WithDeletedAt(deletedAt *time.Time) *PhoneBuilder {
	p.phone.deletedAt = deletedAt
	return p
}

func (p *PhoneBuilder) Build() (*Phone, error) {
	if len(p.errors) > 0 {
		return nil, NewValidationErrors(p.errors)
	}
	if err := p.phone.validate(); err != nil {
		return nil, err
	}
	return p.phone, nil
}

func (p *PhoneBuilder) Apply() error {
	if len(p.errors) > 0 {
		return NewValidationErrors(p.errors)
	}
	if err := p.phone.validate(); err != nil {
		return err
	}
	p.phone.updatedAt = time.Now()
	return nil
}

func (p *Phone) ID() uint               { return p.id }
func (p *Phone) Number() vo.PhoneNumber { return p.number }
func (p *Phone) Kind() vo.Text          { return p.kind }
func (p *Phone) Department() vo.Text    { return p.department }
func (p *Phone) CreatedAt() time.Time   { return p.createdAt }
func (p *Phone) UpdatedAt() time.Time   { return p.updatedAt }
func (p *Phone) DeletedAt() *time.Time  { return p.deletedAt }

func (p *Phone) Delete() error {
	if p.id == 0 {
		return ErrInvalidID
	}
	if p.IsDeleted() {
		return ErrPhoneDeleted
	}
	now := time.Now()
	p.deletedAt = &now
	p.updatedAt = now
	return nil
}

func (p *Phone) IsDeleted() bool {
	return p.deletedAt != nil
}

func (p *Phone) SetID(id uint) error {
	if p.id != 0 {
		return ErrInvalidID
	}
	p.id = id
	return nil
}

func (p *Phone) validate() error {
	return nil
}
