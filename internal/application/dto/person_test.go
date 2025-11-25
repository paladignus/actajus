package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeopleInput(t *testing.T) {
	input := PeopleInput{
		FirstName:  "John",
		LastName:   "Doe",
		BirthDate:  "1990-01-01",
		MotherName: "Jane Doe",
		FatherName: "Jack Doe",
		Gender:     "M",
	}
	t.Run("should return the same values", func(t *testing.T) {
		assert.Equal(t, "John", input.FirstName)
		assert.Equal(t, "Doe", input.LastName)
		assert.Equal(t, "1990-01-01", input.BirthDate)
		assert.Equal(t, "Jane Doe", input.MotherName)
		assert.Equal(t, "Jack Doe", input.FatherName)
		assert.Equal(t, "M", input.Gender)
	})
	t.Run("should return an empty value", func(t *testing.T) {
		input := PeopleInput{}
		assert.Equal(t, "", input.FirstName)
		assert.Equal(t, "", input.LastName)
		assert.Equal(t, "", input.BirthDate)
		assert.Equal(t, "", input.MotherName)
		assert.Equal(t, "", input.FatherName)
		assert.Equal(t, "", input.Gender)
	})
}

