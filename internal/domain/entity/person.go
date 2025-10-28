// Package entity
package entity

import (
	"strings"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

type People struct {
	FirstName vo.Text
	LastName  vo.Text
	BirthDate vo.Date
	Mother    vo.Text
	Father    vo.Text
	Gender    vo.Text
}

func NewPeople(people dto.PeopleInput) (People, error) {
	p := People{
		FirstName: vo.Text(strings.TrimSpace(people.FirstName)),
		LastName:  vo.Text(strings.TrimSpace(people.LastName)),
		BirthDate: vo.Date(strings.TrimSpace(people.BirthDate)),
		Mother:    vo.Text(strings.TrimSpace(people.MotherName)),
		Father:    vo.Text(strings.TrimSpace(people.FatherName)),
		Gender:    vo.Text(strings.TrimSpace(people.Gender)),
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
