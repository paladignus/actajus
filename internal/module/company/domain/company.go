// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// Company representa uma entidade de domínio do contexto de Company.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
//
// Relacionamentos com Address, Phone, Email e SocialMedia são mantidos
// através de IDs de referência, pois estas são entidades compartilhadas
// que podem ser associadas tanto a Company quanto a Person.
type Company struct {
	id               vo.ID
	registeredBy     vo.ID
	idAddress        vo.ID
	idPhone          vo.ID
	idEmail          vo.ID
	idSocialMedia    vo.ID
	registeredByName vo.Text
	name             vo.Text
	tradeName        vo.Text
	cnpj             vo.CNPJ
	createdAt        time.Time
	updatedAt        time.Time
	deletedAt        *time.Time
}

// CompanyBuilder segue o padrão Builder para construção segura de Company.
type CompanyBuilder struct {
	company *Company
}

// NewCompanyBuilder cria um novo builder para Company.
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
	b.company.id = vo.ID(id)
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
	b.company.registeredBy = vo.ID(registeredBy)
	return b
}

func (b *CompanyBuilder) WithRegisteredByName(registeredByName string) *CompanyBuilder {
	b.company.registeredByName = vo.Text(registeredByName)
	return b
}

func (b *CompanyBuilder) WithIDAddress(id int64) *CompanyBuilder {
	b.company.idAddress = vo.ID(id)
	return b
}

func (b *CompanyBuilder) WithCreatedAt(createdAt time.Time) *CompanyBuilder {
	b.company.createdAt = createdAt
	return b
}

func (b *CompanyBuilder) WithIDPhone(id int64) *CompanyBuilder {
	b.company.idPhone = vo.ID(id)
	return b
}

func (b *CompanyBuilder) WithIDEmail(id int64) *CompanyBuilder {
	b.company.idEmail = vo.ID(id)
	return b
}

func (b *CompanyBuilder) WithIDSocialMedia(id int64) *CompanyBuilder {
	b.company.idSocialMedia = vo.ID(id)
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

// Build valida e constrói a entidade Company.
// Retorna erro se os invariantes de domínio forem violados.
func (b *CompanyBuilder) Build() (*Company, error) {
	if err := b.company.validate(); err != nil {
		return nil, err
	}
	return b.company, nil
}

// validate verifica os invariantes de domínio da entidade.
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

// ID retorna o identificador único da empresa.
func (c *Company) ID() vo.ID { return c.id }

// RegisteredBy retorna o identificador de quem registrou a empresa.
func (c *Company) RegisteredBy() vo.ID { return c.registeredBy }

// Name retorna o nome social da empresa.
func (c *Company) Name() vo.Text { return c.name }

// TradeName retorna o nome fantasia da empresa.
func (c *Company) TradeName() vo.Text { return c.tradeName }

// CNPJ retorna o CNPJ da empresa.
func (c *Company) CNPJ() vo.CNPJ { return c.cnpj }

// IDAddress retorna o identificador do endereço principal da empresa.
func (c *Company) IDAddress() vo.ID { return c.idAddress }

// IDPhone retorna o identificador do telefone principal da empresa.
func (c *Company) IDPhone() vo.ID { return c.idPhone }

// IDEmail retorna o identificador do email principal da empresa.
func (c *Company) IDEmail() vo.ID { return c.idEmail }

// IDSocialMedia retorna o identificador da mídia social principal da empresa.
func (c *Company) IDSocialMedia() vo.ID { return c.idSocialMedia }

// RegisteredByName retorna o nome de quem registrou a empresa.
func (c *Company) RegisteredByName() vo.Text { return c.registeredByName }

// CreatedAt retorna a data de criação da empresa.
func (c *Company) CreatedAt() time.Time { return c.createdAt }

// UpdatedAt retorna a data da última atualização da empresa.
func (c *Company) UpdatedAt() time.Time { return c.updatedAt }

// DeletedAt retorna a data de exclusão (nil se não excluída).
func (c *Company) DeletedAt() *time.Time { return c.deletedAt }

// SetID define o ID da empresa apenas uma vez.
func (c *Company) SetID(id int64) error {
	if c.id != 0 {
		return domain.NewFieldError("id", "company ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "company ID is invalid")
	}
	c.id = vo.ID(id)
	return nil
}

// IsDeleted verifica se a empresa foi excluída logicamente.
func (c *Company) IsDeleted() bool {
	return c.deletedAt != nil
}

// HasAddress verifica se a empresa possui um endereço associado.
func (c *Company) HasAddress() bool {
	return c.idAddress != 0
}

// HasPhone verifica se a empresa possui um telefone associado.
func (c *Company) HasPhone() bool {
	return c.idPhone != 0
}

// HasEmail verifica se a empresa possui um email associado.
func (c *Company) HasEmail() bool {
	return c.idEmail != 0
}

// HasSocialMedia verifica se a empresa possui uma mídia social associada.
func (c *Company) HasSocialMedia() bool {
	return c.idSocialMedia != 0
}

// Delete exclui logicamente a empresa.
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

// Restore restaura uma empresa excluída logicamente.
func (c *Company) Restore() {
	if !c.IsDeleted() {
		return
	}
	c.deletedAt = nil
	c.updatedAt = time.Now()
}

// SetAddress associa um endereço à empresa.
func (c *Company) SetAddress(id int64) {
	c.idAddress = vo.ID(id)
	c.updatedAt = time.Now()
}

// SetPhone associa um telefone à empresa.
func (c *Company) SetPhone(id int64) {
	c.idPhone = vo.ID(id)
	c.updatedAt = time.Now()
}

// SetEmail associa um email à empresa.
func (c *Company) SetEmail(id int64) {
	c.idEmail = vo.ID(id)
	c.updatedAt = time.Now()
}

// SetSocialMedia associa uma mídia social à empresa.
func (c *Company) SetSocialMedia(id int64) {
	c.idSocialMedia = vo.ID(id)
	c.updatedAt = time.Now()
}

// UpdateName atualiza o nome social da empresa.
func (c *Company) UpdateName(name string) error {
	if name == "" {
		return domain.NewFieldError("name", "name cannot be empty")
	}
	if !vo.Text(name).IsValid() {
		return domain.NewFieldError("name", "name is invalid")
	}
	c.name = vo.Text(name)
	c.updatedAt = time.Now()
	return nil
}

// UpdateTradeName atualiza o nome fantasia da empresa.
func (c *Company) UpdateTradeName(tradeName string) error {
	if tradeName == "" {
		return domain.NewFieldError("trade_name", "trade name cannot be empty")
	}
	if !vo.Text(tradeName).IsValid() {
		return domain.NewFieldError("trade_name", "trade name is invalid")
	}
	c.tradeName = vo.Text(tradeName)
	c.updatedAt = time.Now()
	return nil
}

// UpdateCNPJ atualiza o CNPJ da empresa.
func (c *Company) UpdateCNPJ(cnpj string) error {
	if cnpj == "" {
		return domain.NewFieldError("cnpj", "CNPJ cannot be empty")
	}
	if !vo.CNPJ(cnpj).IsValid() {
		return domain.NewFieldError("cnpj", "CNPJ is invalid")
	}
	c.cnpj = vo.CNPJ(cnpj)
	c.updatedAt = time.Now()
	return nil
}
