// Package mapper
package mapper

import (
	"strings"

	"github.com/paladignus/actajus/internal/module/person/application/dto"
	"github.com/paladignus/actajus/internal/module/person/domain"
)

type PersonMapper struct{}

func NewPersonMapper() PersonMapper {
	return PersonMapper{}
}

func (p PersonMapper) PersonInputToDomain(input dto.CreatePersonRequest) (*domain.Person, error) {
	name := strings.Fields(strings.TrimSpace(input.Name))
	return domain.NewPersonBuilder().
		WithFirstName(name[0]).
		WithLastName(strings.Join(name[1:], " ")).
		WithBirthday(input.Birthday).
		WithGender(input.Gender).
		Build()
}
