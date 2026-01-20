package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeopleInput(t *testing.T) {
	sut := PersonInput{
		FirstName:  "John",
		LastName:   "Doe",
		BirthDate:  "1990-01-01",
		MotherName: "Jane Doe",
		FatherName: "Jack Doe",
		Gender:     "M",
	}
	t.Run("should return the same values", func(t *testing.T) {
		assert.Equal(t, "John", sut.FirstName)
		assert.Equal(t, "Doe", sut.LastName)
		assert.Equal(t, "1990-01-01", sut.BirthDate)
		assert.Equal(t, "Jane Doe", sut.MotherName)
		assert.Equal(t, "Jack Doe", sut.FatherName)
		assert.Equal(t, "M", sut.Gender)
	})
	t.Run("should return an empty value", func(t *testing.T) {
		sut := PersonInput{}
		assert.Equal(t, "", sut.FirstName)
		assert.Equal(t, "", sut.LastName)
		assert.Equal(t, "", sut.BirthDate)
		assert.Equal(t, "", sut.MotherName)
		assert.Equal(t, "", sut.FatherName)
		assert.Equal(t, "", sut.Gender)
	})
}
