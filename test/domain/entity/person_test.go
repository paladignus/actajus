package entity

import (
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
)

func TestEntityPeople(t *testing.T) {
	_, err := entity.NewPeople(dto.PeopleInput{})
	t.Run("should return error if first name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidFirstName) {
			t.Errorf("Expected error for invalid first name, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: ""})
	t.Run("should return error if last name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidLastName) {
			t.Errorf("Expected error for invalid last name, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: ""})
	t.Run("should return error if birth date is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidBirthDate) {
			t.Errorf("Expected error for invalid birth date, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: ""})
	t.Run("should return error if mother name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidMother) {
			t.Errorf("Expected error for invalid mother name, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: ""})
	t.Run("should return error if father name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidFather) {
			t.Errorf("Expected error for invalid father name, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack"})
	t.Run("should return error if gender is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidGender) {
			t.Errorf("Expected error for invalid gender, but got %v", err)
		}
	})
	_, err = entity.NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack", Gender: "M"})
	t.Run("should return error if marital status is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, entity.ErrInvalidMarital) {
			t.Errorf("Expected error for invalid marital status, but got %v", err)
		}
	})
}
