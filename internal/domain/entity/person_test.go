package entity

import (
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/stretchr/testify/assert"
)

func TestPeople(t *testing.T) {
	_, err := NewPeople(dto.PeopleInput{})
	t.Run("should return error if first name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidFirstName, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: ""})
	t.Run("should return error if last name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidLastName, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: ""})
	t.Run("should return error if birth date is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidBirthDate, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: ""})
	t.Run("should return error if mother name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidMother, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: ""})
	t.Run("should return error if father name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidFather, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack"})
	t.Run("should return error if gender is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidGender, err)
	})
	_, err = NewPeople(dto.PeopleInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack", Gender: "M"})
	t.Run("should null error return if all properties are valid", func(t *testing.T) {
		assert.NoError(t, err)
	})
}
