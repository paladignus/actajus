// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// Phone representa uma entidade de domínio compartilhada entre Company e Person.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
type Phone struct {
	id         vo.ID
	number     vo.PhoneNumber
	kind       vo.Text
	department vo.Text
	createdAt  time.Time
	updatedAt  time.Time
	deletedAt  *time.Time
}

// PhoneBuilder segue o padrão Builder para construção segura de Phone.
type PhoneBuilder struct {
	phone *Phone
}

// NewPhoneBuilder cria um novo builder para Phone.
func NewPhoneBuilder() *PhoneBuilder {
	now := time.Now()
	return &PhoneBuilder{
		phone: &Phone{
			createdAt: now,
			updatedAt: now,
		},
	}
}

func (p *PhoneBuilder) WithID(id int64) *PhoneBuilder {
	p.phone.id = vo.ID(id)
	return p
}

func (p *PhoneBuilder) WithNumber(number string) *PhoneBuilder {
	p.phone.number = vo.PhoneNumber(number)
	return p
}

func (p *PhoneBuilder) WithKind(kind string) *PhoneBuilder {
	p.phone.kind = vo.Text(kind)
	return p
}

func (p *PhoneBuilder) WithDepartment(department string) *PhoneBuilder {
	p.phone.department = vo.Text(department)
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

// Build valida e constrói a entidade Phone.
func (p *PhoneBuilder) Build() (*Phone, error) {
	if err := p.phone.validate(); err != nil {
		return nil, err
	}
	return p.phone, nil
}

// validate verifica os invariantes de domínio da entidade.
func (p *Phone) validate() error {
	if !p.number.IsValid() {
		return domain.NewFieldError("number", "phone number is invalid")
	}
	if !p.kind.IsValid() {
		return domain.NewFieldError("kind", "phone kind is invalid")
	}
	if p.department != "" && !p.department.IsValid() {
		return domain.NewFieldError("department", "phone department is invalid")
	}
	return nil
}

// ID retorna o identificador único do telefone.
func (p *Phone) ID() vo.ID { return p.id }

// Number retorna o número do telefone.
func (p *Phone) Number() vo.PhoneNumber { return p.number }

// Kind retorna o tipo do telefone (ex: "mobile", "commercial").
func (p *Phone) Kind() vo.Text { return p.kind }

// Department retorna o departamento do telefone.
func (p *Phone) Department() vo.Text { return p.department }

// CreatedAt retorna a data de criação do telefone.
func (p *Phone) CreatedAt() time.Time { return p.createdAt }

// UpdatedAt retorna a data da última atualização do telefone.
func (p *Phone) UpdatedAt() time.Time { return p.updatedAt }

// DeletedAt retorna a data de exclusão (nil se não excluído).
func (p *Phone) DeletedAt() *time.Time { return p.deletedAt }

// SetID define o ID do telefone apenas uma vez.
func (p *Phone) SetID(id int64) error {
	if p.id != 0 {
		return domain.NewFieldError("id", "phone ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "phone ID is invalid")
	}
	p.id = vo.ID(id)
	return nil
}

// IsDeleted verifica se o telefone foi excluído logicamente.
func (p *Phone) IsDeleted() bool {
	return p.deletedAt != nil
}

// Delete exclui logicamente o telefone.
func (p *Phone) Delete() error {
	if p.id == 0 {
		return domain.NewFieldError("id", "phone ID is invalid")
	}
	if p.IsDeleted() {
		return domain.NewFieldError("deleted_at", "phone is already deleted")
	}
	now := time.Now()
	p.deletedAt = &now
	p.updatedAt = now
	return nil
}

// Restore restaura um telefone excluído logicamente.
func (p *Phone) Restore() {
	if !p.IsDeleted() {
		return
	}
	p.deletedAt = nil
	p.updatedAt = time.Now()
}

// UpdateNumber atualiza o número do telefone.
func (p *Phone) UpdateNumber(number string) error {
	if !vo.PhoneNumber(number).IsValid() {
		return domain.NewFieldError("number", "phone number is invalid")
	}
	p.number = vo.PhoneNumber(number)
	p.updatedAt = time.Now()
	return nil
}

// UpdateKind atualiza o tipo do telefone.
func (p *Phone) UpdateKind(kind string) error {
	if !vo.Text(kind).IsValid() {
		return domain.NewFieldError("kind", "phone kind is invalid")
	}
	p.kind = vo.Text(kind)
	p.updatedAt = time.Now()
	return nil
}

// UpdateDepartment atualiza o departamento do telefone.
func (p *Phone) UpdateDepartment(department string) {
	p.department = vo.Text(department)
	p.updatedAt = time.Now()
}
