// Package domain
package domain

import (
	"errors"
	"time"
)

type Phone struct {
	id          uint
	contry_code uint
	area_code   uint
	kind        string
	number      string
	department  string
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
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

var (
	ErrInvalidID          = errors.New("invalid phone ID")
	ErrPhoneDeleted       = errors.New("phone is deleted")
	ErrInvalidCountryCode = errors.New("invalid country code")
	ErrInvalidAreaCode    = errors.New("invalid area code")
	ErrInvalidKind        = errors.New("invalid kind")
	ErrInvalidNumber      = errors.New("invalid number")
	ErrInvalidDepartment  = errors.New("invalid department")
)

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

func (p *PhoneBuilder) WithCountryCode(code uint) *PhoneBuilder {
	p.phone.contry_code = code
	if code == 0 {
		p.errors = append(p.errors, ErrInvalidCountryCode)
	}
	return p
}

func (p *PhoneBuilder) WithAreaCode(code uint) *PhoneBuilder {
	p.phone.area_code = code
	if code == 0 {
		p.errors = append(p.errors, ErrInvalidAreaCode)
	}
	return p
}

func (p *PhoneBuilder) WithKind(kind string) *PhoneBuilder {
	p.phone.kind = kind
	if kind == "" {
		p.errors = append(p.errors, ErrInvalidKind)
	}
	return p
}

func (p *PhoneBuilder) WithNumber(number string) *PhoneBuilder {
	p.phone.number = number
	if number == "" {
		p.errors = append(p.errors, ErrInvalidNumber)
	}
	return p
}

func (p *PhoneBuilder) WithDepartment(department string) *PhoneBuilder {
	p.phone.department = department
	if department == "" {
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
		return nil, p.errors
	}
	return p.phone, nil
}

func (p *Phone) SetID(id uint) error {
	if p.id != 0 {
		return ErrInvalidID
	}
	p.id = id
	return nil
}

func (p *Phone) IsDeleted() bool {
	return p.deletedAt != nil
}
