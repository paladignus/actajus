// Package mapper
package mapper

import (
	"github.com/paladignus/actajus/internal/module/person/application/readmodel"
	"github.com/paladignus/actajus/internal/module/person/domain"
)

type PersonProjection struct{}

func NewPersonProjection() PersonProjection {
	return PersonProjection{}
}

func (p PersonProjection) ProjectPersonReadModel(person *domain.Person) *readmodel.PersonReadModel {
	return &readmodel.PersonReadModel{
		ID:        person.ID(),
		Name:      person.FullName(),
		Birthday:  person.Birthday().Value(),
		Gender:    uint32(person.IDGender()),
		CreatedAt: person.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt: person.UpdatedAt().Format("2006-01-02 15:04:05"),
	}
}
