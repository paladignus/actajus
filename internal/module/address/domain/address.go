// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// Address representa uma entidade de domínio compartilhada entre Company e Person.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
//
// Esta entidade é compartilhada, ou seja, tanto Company quanto Person podem
// referenciar o mesmo endereço através de tabelas de relacionamento
// (company_address, person_address).
type Address struct {
	id           vo.ID
	number       uint
	zip          vo.ZIP
	title        vo.Text
	street       vo.Text
	complement   vo.Text
	reference    vo.Text
	neighborhood vo.Text
	city         vo.Text
	state        vo.Text
	country      vo.Text
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

// AddressBuilder segue o padrão Builder para construção segura de Address.
type AddressBuilder struct {
	address *Address
}

// NewAddressBuilder cria um novo builder para Address.
func NewAddressBuilder() *AddressBuilder {
	now := time.Now()
	return &AddressBuilder{
		address: &Address{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (b *AddressBuilder) WithID(id int64) *AddressBuilder {
	b.address.id = vo.ID(id)
	return b
}

func (b *AddressBuilder) WithZIP(zip string) *AddressBuilder {
	b.address.zip = vo.ZIP(zip)
	return b
}

func (b *AddressBuilder) WithTitle(title string) *AddressBuilder {
	b.address.title = vo.Text(title)
	return b
}

func (b *AddressBuilder) WithStreet(street string) *AddressBuilder {
	b.address.street = vo.Text(street)
	return b
}

func (b *AddressBuilder) WithNumber(number uint) *AddressBuilder {
	b.address.number = number
	return b
}

func (b *AddressBuilder) WithComplement(complement string) *AddressBuilder {
	b.address.complement = vo.Text(complement)
	return b
}

func (b *AddressBuilder) WithReference(reference string) *AddressBuilder {
	b.address.reference = vo.Text(reference)
	return b
}

func (b *AddressBuilder) WithNeighborhood(neighborhood string) *AddressBuilder {
	b.address.neighborhood = vo.Text(neighborhood)
	return b
}

func (b *AddressBuilder) WithCity(city string) *AddressBuilder {
	b.address.city = vo.Text(city)
	return b
}

func (b *AddressBuilder) WithState(state string) *AddressBuilder {
	b.address.state = vo.Text(state)
	return b
}

func (b *AddressBuilder) WithCountry(country string) *AddressBuilder {
	b.address.country = vo.Text(country)
	return b
}

func (b *AddressBuilder) WithCreatedAt(createdAt time.Time) *AddressBuilder {
	b.address.createdAt = createdAt
	return b
}

func (b *AddressBuilder) WithUpdatedAt(updatedAt time.Time) *AddressBuilder {
	b.address.updatedAt = updatedAt
	return b
}

func (b *AddressBuilder) WithDeletedAt(deletedAt *time.Time) *AddressBuilder {
	b.address.deletedAt = deletedAt
	return b
}

// Build valida e constrói a entidade Address.
// Retorna erro se os invariantes de domínio forem violados.
func (b *AddressBuilder) Build() (*Address, error) {
	if err := b.address.validate(); err != nil {
		return nil, err
	}
	return b.address, nil
}

// validate verifica os invariantes de domínio da entidade.
func (a *Address) validate() error {
	if !a.zip.IsValid() {
		return domain.NewFieldError("zip", "address ZIP is invalid")
	}
	if !a.title.IsValid() {
		return domain.NewFieldError("title", "address title is invalid")
	}
	if !a.street.IsValid() {
		return domain.NewFieldError("street", "address street is invalid")
	}
	if a.number == 0 {
		return domain.NewFieldError("number", "address number is invalid")
	}
	if !a.neighborhood.IsValid() {
		return domain.NewFieldError("neighborhood", "address neighborhood is invalid")
	}
	if !a.city.IsValid() {
		return domain.NewFieldError("city", "address city is invalid")
	}
	if !a.state.IsValid() {
		return domain.NewFieldError("state", "address state is invalid")
	}
	if !a.country.IsValid() {
		return domain.NewFieldError("country", "address country is invalid")
	}
	return nil
}

// ID retorna o identificador único do endereço.
func (a *Address) ID() vo.ID { return a.id }

// Number retorna o número do endereço.
func (a *Address) Number() uint { return a.number }

// ZIP retorna o CEP do endereço.
func (a *Address) ZIP() vo.ZIP { return a.zip }

// Title retorna o título do endereço (ex: "Casa", "Trabalho").
func (a *Address) Title() vo.Text { return a.title }

// Street retorna a rua do endereço.
func (a *Address) Street() vo.Text { return a.street }

// Complement retorna o complemento do endereço.
func (a *Address) Complement() vo.Text { return a.complement }

// Reference retorna o ponto de referência do endereço.
func (a *Address) Reference() vo.Text { return a.reference }

// Neighborhood retorna o bairro do endereço.
func (a *Address) Neighborhood() vo.Text { return a.neighborhood }

// City retorna a cidade do endereço.
func (a *Address) City() vo.Text { return a.city }

// State retorna o estado do endereço.
func (a *Address) State() vo.Text { return a.state }

// Country retorna o país do endereço.
func (a *Address) Country() vo.Text { return a.country }

// CreatedAt retorna a data de criação do endereço.
func (a *Address) CreatedAt() time.Time { return a.createdAt }

// UpdatedAt retorna a data da última atualização do endereço.
func (a *Address) UpdatedAt() time.Time { return a.updatedAt }

// DeletedAt retorna a data de exclusão (nil se não excluído).
func (a *Address) DeletedAt() *time.Time { return a.deletedAt }

// SetID define o ID do endereço apenas uma vez.
func (a *Address) SetID(id int64) error {
	if a.id != 0 {
		return domain.NewFieldError("id", "address ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "address ID is invalid")
	}
	a.id = vo.ID(id)
	return nil
}

// IsDeleted verifica se o endereço foi excluído logicamente.
func (a *Address) IsDeleted() bool {
	return a.deletedAt != nil
}

// Delete exclui logicamente o endereço.
func (a *Address) Delete() error {
	if a.id == 0 {
		return domain.NewFieldError("id", "address ID is invalid")
	}
	if a.IsDeleted() {
		return domain.NewFieldError("deleted_at", "address is already deleted")
	}
	now := time.Now()
	a.deletedAt = &now
	a.updatedAt = now
	return nil
}

// Restore restaura um endereço excluído logicamente.
func (a *Address) Restore() {
	if !a.IsDeleted() {
		return
	}
	a.deletedAt = nil
	a.updatedAt = time.Now()
}

// UpdateTitle atualiza o título do endereço.
func (a *Address) UpdateTitle(title string) error {
	if title == "" {
		return domain.NewFieldError("title", "title cannot be empty")
	}
	if !vo.Text(title).IsValid() {
		return domain.NewFieldError("title", "title is invalid")
	}
	a.title = vo.Text(title)
	a.updatedAt = time.Now()
	return nil
}

// UpdateStreet atualiza a rua do endereço.
func (a *Address) UpdateStreet(street string) error {
	if street == "" {
		return domain.NewFieldError("street", "street cannot be empty")
	}
	if !vo.Text(street).IsValid() {
		return domain.NewFieldError("street", "street is invalid")
	}
	a.street = vo.Text(street)
	a.updatedAt = time.Now()
	return nil
}

// UpdateNumber atualiza o número do endereço.
func (a *Address) UpdateNumber(number uint) {
	a.number = number
	a.updatedAt = time.Now()
}

// UpdateComplement atualiza o complemento do endereço.
func (a *Address) UpdateComplement(complement string) {
	a.complement = vo.Text(complement)
	a.updatedAt = time.Now()
}

// UpdateReference atualiza o ponto de referência do endereço.
func (a *Address) UpdateReference(reference string) {
	a.reference = vo.Text(reference)
	a.updatedAt = time.Now()
}

// UpdateNeighborhood atualiza o bairro do endereço.
func (a *Address) UpdateNeighborhood(neighborhood string) error {
	if neighborhood == "" {
		return domain.NewFieldError("neighborhood", "neighborhood cannot be empty")
	}
	if !vo.Text(neighborhood).IsValid() {
		return domain.NewFieldError("neighborhood", "neighborhood is invalid")
	}
	a.neighborhood = vo.Text(neighborhood)
	a.updatedAt = time.Now()
	return nil
}

// UpdateCity atualiza a cidade do endereço.
func (a *Address) UpdateCity(city string) error {
	if city == "" {
		return domain.NewFieldError("city", "city cannot be empty")
	}
	if !vo.Text(city).IsValid() {
		return domain.NewFieldError("city", "city is invalid")
	}
	a.city = vo.Text(city)
	a.updatedAt = time.Now()
	return nil
}

// UpdateState atualiza o estado do endereço.
func (a *Address) UpdateState(state string) error {
	if state == "" {
		return domain.NewFieldError("state", "state cannot be empty")
	}
	if !vo.Text(state).IsValid() {
		return domain.NewFieldError("state", "state is invalid")
	}
	a.state = vo.Text(state)
	a.updatedAt = time.Now()
	return nil
}

// UpdateCountry atualiza o país do endereço.
func (a *Address) UpdateCountry(country string) error {
	if country == "" {
		return domain.NewFieldError("country", "country cannot be empty")
	}
	if !vo.Text(country).IsValid() {
		return domain.NewFieldError("country", "country is invalid")
	}
	a.country = vo.Text(country)
	a.updatedAt = time.Now()
	return nil
}

// UpdateZIP atualiza o CEP do endereço.
func (a *Address) UpdateZIP(zip string) error {
	if !vo.ZIP(zip).IsValid() {
		return domain.NewFieldError("zip", "ZIP is invalid")
	}
	a.zip = vo.ZIP(zip)
	a.updatedAt = time.Now()
	return nil
}

// FullAddress retorna o endereço completo formatado.
// Ex: "Rua das Flores, 123, Centro, São Paulo - SP, 01000-000"
func (a *Address) FullAddress() string {
	var parts []string

	if a.street != "" {
		parts = append(parts, a.street.Value())
	}
	if a.number > 0 {
		parts = append(parts, string(rune(a.number+'0')))
	}
	if a.complement != "" {
		parts = append(parts, a.complement.Value())
	}
	if a.neighborhood != "" {
		parts = append(parts, a.neighborhood.Value())
	}
	if a.city != "" {
		parts = append(parts, a.city.Value())
	}
	if a.state != "" {
		parts = append(parts, a.state.Value())
	}
	if a.zip != "" {
		parts = append(parts, a.zip.Value())
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}
	return result
}
