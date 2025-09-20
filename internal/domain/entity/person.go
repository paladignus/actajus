// Package entity
package entity

import (
	"errors"
	"strings"

	"github.com/paladignus/actajus/internal/application/dto"
	vo "github.com/paladignus/actajus/internal/domain/value_object"
)

var (
	ErrInvalidFirstName = errors.New("first name is invalid")
	ErrInvalidLastName  = errors.New("last name is invalid")
	ErrInvalidBirthDate = errors.New("birth date is invalid")
	ErrInvalidMother    = errors.New("mother is invalid")
	ErrInvalidFather    = errors.New("father is invalid")
	ErrInvalidGender    = errors.New("gender is invalid")
	ErrInvalidMarital   = errors.New("marital status is invalid")
)

type People struct {
	FirstName     vo.Text
	LastName      vo.Text
	BirthDate     vo.Date
	Mother        vo.Text
	Father        vo.Text
	Gender        vo.Text
	MaritalStatus vo.Text
}

func NewPeople(people dto.PeopleInputDTO) (People, error) {
	p := People{
		FirstName:     vo.Text(strings.TrimSpace(people.FirstName)),
		LastName:      vo.Text(strings.TrimSpace(people.LastName)),
		BirthDate:     vo.Date(strings.TrimSpace(people.BirthDate)),
		Mother:        vo.Text(strings.TrimSpace(people.MotherName)),
		Father:        vo.Text(strings.TrimSpace(people.FatherName)),
		Gender:        vo.Text(strings.TrimSpace(people.Gender)),
		MaritalStatus: vo.Text(strings.TrimSpace(people.MaritalStatus)),
	}
	if !p.FirstName.IsValid() {
		return p, ErrInvalidFirstName
	}
	if !p.LastName.IsValid() {
		return p, ErrInvalidLastName
	}
	if !p.BirthDate.IsValid() {
		return p, ErrInvalidBirthDate
	}
	if !p.Mother.IsValid() {
		return p, ErrInvalidMother
	}
	if !p.Father.IsValid() {
		return p, ErrInvalidFather
	}
	if !p.Gender.IsValid() {
		return p, ErrInvalidGender
	}
	if !p.MaritalStatus.IsValid() {
		return p, ErrInvalidMarital
	}
	return p, nil
}
