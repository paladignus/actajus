package entity

import (
	"errors"
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
)

func TestPeople(t *testing.T) {
	_, err := NewPeople(dto.PeopleInput{})
	t.Run("should return error if first name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidFirstName) {
			t.Errorf("Expected error for invalid first name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: ""})
	t.Run("should return error if last name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidLastName) {
			t.Errorf("Expected error for invalid last name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: ""})
	t.Run("should return error if birth date is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidBirthDate) {
			t.Errorf("Expected error for invalid birth date, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: ""})
	t.Run("should return error if mother name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidMother) {
			t.Errorf("Expected error for invalid mother name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: ""})
	t.Run("should return error if father name is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidFather) {
			t.Errorf("Expected error for invalid father name, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack"})
	t.Run("should return error if gender is invalid", func(t *testing.T) {
		if err == nil || !errors.Is(err, ErrInvalidGender) {
			t.Errorf("Expected error for invalid gender, but got %v", err)
		}
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack", Gender: "M"})
	t.Run("should null error return if all properties are valid", func(t *testing.T) {
		if err != nil {
			t.Errorf("Expected null error return if all properties are valid, but got %v", err)
		}
	})
}
