// Package entity provides the entity of the application
package entity

import (
	"errors"
	"strings"
	"testing"

	"github.com/paladignus/juridico/internal/application/dto"
	vo "github.com/paladignus/juridico/internal/domain/value_object"
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

func TestEntityPeople(t *testing.T) {
	_, err := NewPeople(dto.PeopleInputDTO{})
	t.Run("should return error if first name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidFirstName) {
			t.Errorf("Expected error for invalid first name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInputDTO{FirstName: "John", LastName: ""})
	t.Run("should return error if last name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidLastName) {
			t.Errorf("Expected error for invalid last name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInputDTO{FirstName: "John", LastName: "Doe", BirthDate: ""})
	t.Run("should return error if birth date is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidBirthDate) {
			t.Errorf("Expected error for invalid birth date, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInputDTO{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: ""})
	t.Run("should return error if mother name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidMother) {
			t.Errorf("Expected error for invalid mother name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInputDTO{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: ""})
	t.Run("should return error if father name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidFather) {
			t.Errorf("Expected error for invalid father name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInputDTO{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack"})
	t.Run("should return error if gender is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidGender) {
			t.Errorf("Expected error for invalid gender, but got %v", err)
		}
	})
}
