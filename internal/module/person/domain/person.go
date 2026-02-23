// Package domain
package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/paladignus/actajus/internal/shared/domain"
	vo "github.com/paladignus/actajus/internal/shared/domain/value_object"
)

type Person struct {
	id        uint
	firstName vo.Text
	lastName  vo.Text
	idGender  uint
	birthday  vo.Date
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

type PersonBuilder struct {
	person *Person
	errors []domain.Violation
}

func NewPersonBuilder() *PersonBuilder {
	now := time.Now()
	return &PersonBuilder{
		person: &Person{
			createdAt: now,
			updatedAt: now,
		},
		errors: make([]domain.Violation, 0),
	}
}

func (p *PersonBuilder) WithID(id uint) *PersonBuilder {
	p.person.id = id
	return p
}

func (p *PersonBuilder) WithFirstName(firstName string) *PersonBuilder {
	p.person.firstName = vo.Text(firstName)
	return p
}

func (p *PersonBuilder) WithLastName(lastName string) *PersonBuilder {
	p.person.lastName = vo.Text(lastName)
	return p
}

func (p *PersonBuilder) WithGender(idGender uint) *PersonBuilder {
	p.person.idGender = idGender
	return p
}

func (p *PersonBuilder) WithBirthday(birthday string) *PersonBuilder {
	p.person.birthday = vo.Date(birthday)
	return p
}

func (p *PersonBuilder) WithCreatedAt(createdAt time.Time) *PersonBuilder {
	p.person.createdAt = createdAt
	return p
}

func (p *PersonBuilder) WithUpdatedAt(updatedAt time.Time) *PersonBuilder {
	p.person.updatedAt = updatedAt
	return p
}

func (p *PersonBuilder) WithDeletedAt(deletedAt *time.Time) *PersonBuilder {
	p.person.deletedAt = deletedAt
	return p
}

func (p PersonBuilder) Build() (*Person, error) {
	p.validate()
	if len(p.errors) > 0 {
		return nil, domain.NewValidationError(p.errors)
	}
	return p.person, nil
}

func (p *Person) ID() uint              { return p.id }
func (p *Person) FirstName() vo.Text    { return p.firstName }
func (p *Person) LastName() vo.Text     { return p.lastName }
func (p *Person) FullName() string      { return p.firstName.Value() + " " + p.lastName.Value() }
func (p *Person) IDGender() uint        { return p.idGender }
func (p *Person) Birthday() vo.Date     { return p.birthday }
func (p *Person) CreatedAt() time.Time  { return p.createdAt }
func (p *Person) UpdatedAt() time.Time  { return p.updatedAt }
func (p *Person) DeletedAt() *time.Time { return p.deletedAt }

func (p *Person) SetID(id uint) error {
	if p.id != 0 {
		return errors.New("id is already set")
	}
	if id == 0 {
		return errors.New("id must be greater than 0")
	}
	p.id = id
	return nil
}

func (p *Person) Delete() error {
	if p.id == 0 {
		return errors.New("person not found")
	}
	if p.IsDeleted() {
		return errors.New("person already deleted")
	}
	now := time.Now()
	p.deletedAt = &now
	return nil
}

func (p *Person) IsDeleted() bool {
	return p.deletedAt != nil
}

func (p *PersonBuilder) validate() {
	if strings.TrimSpace(p.person.firstName.Value()) == "" {
		p.errors = append(p.errors, domain.Violation{
			Path: "first_name",
			Code: domain.CodeRequired,
		})
	} else if !p.person.firstName.IsValid() {
		p.errors = append(p.errors, domain.Violation{
			Path: "first_name",
			Code: domain.CodeInvalid,
			Meta: map[string]string{"type": "text"},
		})
	}
	if strings.TrimSpace(p.person.lastName.Value()) == "" {
		p.errors = append(p.errors, domain.Violation{
			Path: "last_name",
			Code: domain.CodeRequired,
		})
	} else if !p.person.lastName.IsValid() {
		p.errors = append(p.errors, domain.Violation{
			Path: "last_name",
			Code: domain.CodeInvalid,
			Meta: map[string]string{"type": "text"},
		})
	}
	if p.person.idGender == 0 {
		p.errors = append(p.errors, domain.Violation{
			Path: "id_gender",
			Code: domain.CodeRequired,
		})
	}
	if strings.TrimSpace(p.person.birthday.Value()) == "" {
		p.errors = append(p.errors, domain.Violation{
			Path: "birthday",
			Code: domain.CodeRequired,
		})
	} else if !p.person.birthday.IsValid() {
		p.errors = append(p.errors, domain.Violation{
			Path: "birthday",
			Code: domain.CodeInvalid,
			Meta: map[string]string{"format": "dd/mm/yyyy"},
		})
	}
}
