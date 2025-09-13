// Package entity provides the entity of the application
package entity

import (
	"testing"

	"github.com/paladignus/juridico/internal/application/dto"
)

type People struct {
	FirstName     string
	LastName      string
	BirthDate     string
	MotherName    string
	FatherName    string
	Gender        string
	MaritalStatus string
}

func NewPeople(people dto.PeopleInputDTO) People {
	return People{
		FirstName:     people.FirstName,
		LastName:      people.LastName,
		BirthDate:     people.BirthDate,
		MotherName:    people.MotherName,
		FatherName:    people.FatherName,
		Gender:        people.Gender,
		MaritalStatus: people.MaritalStatus,
	}
}

func (p People) IsValid() error {
}

func TestEntityPeople(t *testing.T) {
	sut := NewPeople(dto.PeopleInputDTO{
		FirstName:     "John",
		LastName:      "Doe",
		BirthDate:     "1980-02-30",
		MotherName:    "Helena Doe",
		FatherName:    "Martin Doe",
		Gender:        "Masculino",
		MaritalStatus: "Solteiro",
	})
	t.Run("should return error if FirstName is empty", func(t *testing.T) {
		sut.FirstName = ""
		if err := sut.IsValid(); err == nil {
			t.Errorf("Expected error for empty name, but got nil")
		}
	})
}
