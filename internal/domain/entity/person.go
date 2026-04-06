// Package entity
package entity

import (
	"strings"
	"time"

	"github.com/paladignus/actajus/internal/application/command"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type Person struct {
	IDPerson  int
	FirstName vo.Text
	LastName  vo.Text
	BirthDate vo.Date
	Mother    vo.Text
	Father    vo.Text
	Gender    vo.Text
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewPerson(input command.CreatePersonCommand) (Person, error) {
	p := Person{
		FirstName: vo.Text(strings.TrimSpace(input.FirstName)),
		LastName:  vo.Text(strings.TrimSpace(input.LastName)),
		BirthDate: vo.Date(strings.TrimSpace(input.BirthDate)),
		Mother:    vo.Text(strings.TrimSpace(input.MotherName)),
		Father:    vo.Text(strings.TrimSpace(input.FatherName)),
		Gender:    vo.Text(strings.TrimSpace(input.Gender)),
	}
	if !p.FirstName.IsValid() {
		return p, exception.ErrInvalidFirstName
	}
	if !p.LastName.IsValid() {
		return p, exception.ErrInvalidLastName
	}
	if !p.BirthDate.IsValid() {
		return p, exception.ErrInvalidBirthDate
	}
	if !p.Mother.IsValid() {
		return p, exception.ErrInvalidMother
	}
	if !p.Father.IsValid() {
		return p, exception.ErrInvalidFather
	}
	if !p.Gender.IsValid() {
		return p, exception.ErrInvalidGender
	}
	return p, nil
}
