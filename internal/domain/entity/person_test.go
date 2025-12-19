// Package entity
package entity

import (
	"testing"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/exception"
	"github.com/stretchr/testify/assert"
)

func TestPerson(t *testing.T) {
	_, err := NewPerson(dto.PersonInput{})
	t.Run("should return error if first name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidFirstName, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: ""})
	t.Run("should return error if last name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidLastName, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: "Doe", BirthDate: ""})
	t.Run("should return error if birth date is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidBirthDate, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: ""})
	t.Run("should return error if mother name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidMother, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: ""})
	t.Run("should return error if father name is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidFather, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack"})
	t.Run("should return error if gender is invalid", func(t *testing.T) {
		assert.Error(t, err)
		assert.Equal(t, exception.ErrInvalidGender, err)
	})
	_, err = NewPerson(dto.PersonInput{FirstName: "John", LastName: "Doe", BirthDate: "28/08/1999", MotherName: "Jane", FatherName: "Jack", Gender: "M"})
	t.Run("should null error return if all properties are valid", func(t *testing.T) {
		assert.NoError(t, err)
	})
}
