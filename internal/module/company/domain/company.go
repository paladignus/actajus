// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Company struct {
	id           uint
	registeredBy uint
	idAddress    uint
	idPhone      uint
	name         vo.Text
	tradeName    vo.Text
	cnpj         vo.CNPJ
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

type CompanyBuilder struct {
	company *Company
}

func NewCompanyBuilder() *CompanyBuilder {
	now := time.Now()
	return &CompanyBuilder{
		company: &Company{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (c *Company) UpdateBuilder() *CompanyBuilder {
	return &CompanyBuilder{
		company: c,
	}
}

func (b *CompanyBuilder) WithID(id uint) *CompanyBuilder {
	b.company.id = id
	return b
}

func (b *CompanyBuilder) WithName(name string) *CompanyBuilder {
	b.company.name = vo.Text(name)
	return b
}

func (b *CompanyBuilder) WithTradeName(tradeName string) *CompanyBuilder {
	b.company.tradeName = vo.Text(tradeName)
	return b
}

func (b *CompanyBuilder) WithCNPJ(cnpj string) *CompanyBuilder {
	b.company.cnpj = vo.CNPJ(cnpj)
	return b
}

func (b *CompanyBuilder) WithRegisteredBy(registeredBy uint) *CompanyBuilder {
	b.company.registeredBy = registeredBy
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
	if err := b.company.validate(); err != nil {
		return nil, b.company.validate()
	}
	return b.company, nil
}

func (b *CompanyBuilder) Apply() error {
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
func (c *Company) IDPhone() uint         { return c.idPhone }
func (c *Company) CreatedAt() time.Time  { return c.createdAt }
func (c *Company) UpdatedAt() time.Time  { return c.updatedAt }
func (c *Company) DeletedAt() *time.Time { return c.deletedAt }

func (c *Company) Delete() error {
	if c.id == 0 {
		return domain.NewFieldError("id", "company ID is invalid")
	}
	if c.IsDeleted() {
		return domain.NewFieldError("deleted_at", "company is already deleted")
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

func (c *Company) SetPhone(id uint) {
	c.idPhone = id
	c.updatedAt = time.Now()
}

func (c *Company) HasAddress() bool {
	return c.idAddress != 0
}

func (c *Company) SetID(id uint) error {
	if c.id != 0 {
		return domain.NewFieldError("id", "company ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "company ID is invalid")
	}
	c.id = id
	return nil
}

func (c *Company) validate() error {
	if !c.name.IsValid() {
		return domain.NewFieldError("name", "name is invalid")
	}
	if !c.tradeName.IsValid() {
		return domain.NewFieldError("trade_name", "trade name is invalid")
	}
	if !c.cnpj.IsValid() {
		return domain.NewFieldError("cnpj", "CNPJ check digits are invalid")
	}
	if c.registeredBy == 0 {
		return domain.NewFieldError("registered_by", "registered by is invalid")
	}
	return nil
}
