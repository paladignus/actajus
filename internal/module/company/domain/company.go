// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Company struct {
	id               int64
	registeredBy     int64
	idAddress        int64
	idPhone          int64
	idEmail          int64
	registeredByName vo.Text
	name             vo.Text
	tradeName        vo.Text
	cnpj             vo.CNPJ
	createdAt        time.Time
	updatedAt        time.Time
	deletedAt        *time.Time
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

func (b *CompanyBuilder) WithID(id int64) *CompanyBuilder {
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

func (b *CompanyBuilder) WithRegisteredBy(registeredBy int64) *CompanyBuilder {
	b.company.registeredBy = registeredBy
	return b
}

func (b *CompanyBuilder) WithRegisteredByName(registeredByName string) *CompanyBuilder {
	b.company.registeredByName = vo.Text(registeredByName)
	return b
}

func (b *CompanyBuilder) WithIDAddress(id int64) *CompanyBuilder {
	b.company.idAddress = id
	return b
}

func (b *CompanyBuilder) WithCreatedAt(createdAt time.Time) *CompanyBuilder {
	b.company.createdAt = createdAt
	return b
}

func (b *CompanyBuilder) WithIDPhone(id int64) *CompanyBuilder {
	b.company.idPhone = id
	return b
}

func (b *CompanyBuilder) WithIDEmail(id int64) *CompanyBuilder {
	b.company.idEmail = id
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

func (c *Company) ID() int64                 { return c.id }
func (c *Company) RegisteredBy() int64       { return c.registeredBy }
func (c *Company) Name() vo.Text             { return c.name }
func (c *Company) TradeName() vo.Text        { return c.tradeName }
func (c *Company) CNPJ() vo.CNPJ             { return c.cnpj }
func (c *Company) IDAddress() int64          { return c.idAddress }
func (c *Company) IDPhone() int64            { return c.idPhone }
func (c *Company) IDEmail() int64            { return c.idEmail }
func (c *Company) RegisteredByName() vo.Text { return c.registeredByName }
func (c *Company) CreatedAt() time.Time      { return c.createdAt }
func (c *Company) UpdatedAt() time.Time      { return c.updatedAt }
func (c *Company) DeletedAt() *time.Time     { return c.deletedAt }

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

func (c *Company) SetAddress(id int64) {
	c.idAddress = id
	c.updatedAt = time.Now()
}

func (c *Company) SetPhone(id int64) {
	c.idPhone = id
	c.updatedAt = time.Now()
}

func (c *Company) SetEmail(id int64) {
	c.idEmail = id
	c.updatedAt = time.Now()
}

func (c *Company) HasAddress() bool {
	return c.idAddress != 0
}

func (c *Company) SetID(id int64) error {
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
	if c.name != "" && !c.name.IsValid() {
		return domain.NewFieldError("name", "name is invalid")
	}
	if c.tradeName != "" && !c.tradeName.IsValid() {
		return domain.NewFieldError("trade_name", "trade name is invalid")
	}
	if c.cnpj != "" && !c.cnpj.IsValid() {
		return domain.NewFieldError("cnpj", "CNPJ check digits are invalid")
	}
	if c.id == 0 && c.registeredBy == 0 {
		return domain.NewFieldError("registered_by", "registered by is invalid")
	}
	return nil
}
