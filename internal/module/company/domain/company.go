// Package domain
package domain

import (
	"time"

	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Company struct {
	id           uint
	registeredBy uint
	idAddress    uint
	name         vo.Text
	tradeName    vo.Text
	cnpj         vo.CNPJ
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

type CompanyBuilder struct {
	company *Company
	errors  []error
}

func NewCompanyBuilder() *CompanyBuilder {
	now := time.Now()
	return &CompanyBuilder{
		company: &Company{
			createdAt: now,
			updatedAt: now,
		},
		errors: []error{},
	}
}

func (c *Company) UpdateBuilder() *CompanyBuilder {
	if c.id == 0 {
		return &CompanyBuilder{
			errors: []error{ErrInvalidID},
		}
	}
	if c.IsDeleted() {
		return &CompanyBuilder{
			errors: []error{ErrCompanyDeleted},
		}
	}
	return &CompanyBuilder{
		company: c,
		errors:  []error{},
	}
}

func (b *CompanyBuilder) WithID(id uint) *CompanyBuilder {
	b.company.id = id
	if id == 0 {
		b.errors = append(b.errors, ErrInvalidID)
	}
	return b
}

func (b *CompanyBuilder) WithName(name string) *CompanyBuilder {
	b.company.name = vo.Text(name)
	if !b.company.name.IsValid() {
		b.errors = append(b.errors, ErrInvalidName)
	}
	return b
}

func (b *CompanyBuilder) WithTradeName(tradeName string) *CompanyBuilder {
	b.company.tradeName = vo.Text(tradeName)
	if !b.company.tradeName.IsValid() {
		b.errors = append(b.errors, ErrInvalidTradeName)
	}
	return b
}

func (b *CompanyBuilder) WithCNPJ(cnpj string) *CompanyBuilder {
	b.company.cnpj = vo.CNPJ(cnpj)
	if !b.company.cnpj.IsValid() {
		b.errors = append(b.errors, ErrInvalidCNPJ)
	}
	return b
}

func (b *CompanyBuilder) WithRegisteredBy(registeredBy uint) *CompanyBuilder {
	b.company.registeredBy = registeredBy
	if registeredBy <= 0 {
		b.errors = append(b.errors, ErrInvalidRegisteredBy)
	}
	return b
}

func (b *CompanyBuilder) WithIDAddress(id uint) *CompanyBuilder {
	b.company.idAddress = id
	return b
}

func (b *CompanyBuilder) WithCreatedAt(createdAt time.Time) *CompanyBuilder {
	b.company.createdAt = createdAt
	return b
}

func (b *CompanyBuilder) WithUpdatedAt(updatedAt time.Time) *CompanyBuilder {
	b.company.updatedAt = updatedAt
	return b
}

func (b *CompanyBuilder) WithDeletedAt(deletedAt *time.Time) *CompanyBuilder {
	b.company.deletedAt = deletedAt
	return b
}

func (b *CompanyBuilder) Build() (*Company, error) {
	if len(b.errors) > 0 {
		return nil, NewValidationErrors(b.errors)
	}
	if err := b.company.validate(); err != nil {
		return nil, err
	}
	return b.company, nil
}

func (b *CompanyBuilder) Apply() error {
	if len(b.errors) > 0 {
		return NewValidationErrors(b.errors)
	}
	if err := b.company.validate(); err != nil {
		return err
	}
	b.company.updatedAt = time.Now()
	return nil
}

func (c *Company) ID() uint              { return c.id }
func (c *Company) RegisteredBy() uint    { return c.registeredBy }
func (c *Company) Name() vo.Text         { return c.name }
func (c *Company) TradeName() vo.Text    { return c.tradeName }
func (c *Company) CNPJ() vo.CNPJ         { return c.cnpj }
func (c *Company) IDAddress() uint       { return c.idAddress }
func (c *Company) CreatedAt() time.Time  { return c.createdAt }
func (c *Company) UpdatedAt() time.Time  { return c.updatedAt }
func (c *Company) DeletedAt() *time.Time { return c.deletedAt }

func (c *Company) Delete() error {
	if c.id == 0 {
		return ErrInvalidID
	}
	if c.IsDeleted() {
		return ErrAlreadyDeleted
	}
	now := time.Now()
	c.deletedAt = &now
	c.updatedAt = now
	return nil
}

func (c *Company) IsDeleted() bool {
	return c.deletedAt != nil
}

func (c *Company) SetAddress(id uint) {
	c.idAddress = id
	c.updatedAt = time.Now()
}

func (c *Company) HasAddress() bool {
	return c.idAddress != 0
}

func (c *Company) SetID(id uint) error {
	if c.id != 0 {
		return ErrIDAlreadySet
	}
	if id == 0 {
		return ErrInvalidID
	}
	c.id = id
	return nil
}

func (c *Company) validate() error {
	if !c.name.IsValid() {
		return ErrInvalidName
	}
	if !c.tradeName.IsValid() {
		return ErrInvalidTradeName
	}
	if !c.cnpj.IsValid() {
		return ErrInvalidCNPJ
	}
	if c.registeredBy <= 0 {
		return ErrInvalidRegisteredBy
	}
	return nil
}
