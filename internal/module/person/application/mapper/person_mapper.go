// Package mapper
package mapper

import (
	"strings"

	"github.com/paladignus/actajus/internal/module/person/application/command"
	"github.com/paladignus/actajus/internal/module/person/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type PersonMapper struct{}

func NewPersonMapper() PersonMapper {
	return PersonMapper{}
}

func (p PersonMapper) PersonInputToDomain(input command.CreatePersonCommand) (*domain.Person, error) {
	name := strings.Fields(strings.TrimSpace(input.Name))
	if len(name) == 0 {
		return nil, sharedDomain.NewValidationError([]sharedDomain.Violation{{
			Path: "name",
			Code: sharedDomain.CodeRequired,
		}})
	}
	if len(name) < 2 {
		return nil, sharedDomain.NewValidationError([]sharedDomain.Violation{{
			Path: "name",
			Code: sharedDomain.CodeInvalid,
			Meta: map[string]string{"reason": "full_name_required"},
		}})
	}
	return domain.NewPersonBuilder().
		WithFirstName(name[0]).
		WithLastName(strings.Join(name[1:], " ")).
		WithBirthday(input.Birthday).
		WithGender(input.Gender).
		Build()
}
