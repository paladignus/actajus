// Package domain
package domain

import (
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

// Email representa uma entidade de domínio compartilhada entre Company e Person.
// Como Entidade Rica, ela contém lógica de negócio e protege seus invariantes.
type Email struct {
	id        vo.ID
	address   vo.Email
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

// EmailBuilder segue o padrão Builder para construção segura de Email.
type EmailBuilder struct {
	email *Email
}

// NewEmailBuilder cria um novo builder para Email.
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
	e.email.id = vo.ID(id)
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

// Build valida e constrói a entidade Email.
func (e *EmailBuilder) Build() (*Email, error) {
	if err := e.email.validate(); err != nil {
		return nil, err
	}
	return e.email, nil
}

// validate verifica os invariantes de domínio da entidade.
func (e *Email) validate() error {
	if !e.address.IsValid() {
		return domain.NewFieldError("address", "email address is invalid")
	}
	return nil
}

// ID retorna o identificador único do email.
func (e *Email) ID() vo.ID { return e.id }

// Address retorna o endereço de email.
func (e *Email) Address() vo.Email { return e.address }

// CreatedAt retorna a data de criação do email.
func (e *Email) CreatedAt() time.Time { return e.createdAt }

// UpdatedAt retorna a data da última atualização do email.
func (e *Email) UpdatedAt() time.Time { return e.updatedAt }

// DeletedAt retorna a data de exclusão (nil se não excluído).
func (e *Email) DeletedAt() *time.Time { return e.deletedAt }

// SetID define o ID do email apenas uma vez.
func (e *Email) SetID(id int64) error {
	if e.id != 0 {
		return domain.NewFieldError("id", "email ID is already set")
	}
	if id == 0 {
		return domain.NewFieldError("id", "email ID is invalid")
	}
	e.id = vo.ID(id)
	return nil
}

// IsDeleted verifica se o email foi excluído logicamente.
func (e *Email) IsDeleted() bool {
	return e.deletedAt != nil
}

// Delete exclui logicamente o email.
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

// Restore restaura um email excluído logicamente.
func (e *Email) Restore() {
	if !e.IsDeleted() {
		return
	}
	e.deletedAt = nil
	e.updatedAt = time.Now()
}

// UpdateAddress atualiza o endereço de email.
func (e *Email) UpdateAddress(address string) error {
	if !vo.Email(address).IsValid() {
		return domain.NewFieldError("address", "email address is invalid")
	}
	e.address = vo.Email(address)
	e.updatedAt = time.Now()
	return nil
}
