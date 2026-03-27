package mapper_test

import (
	"testing"

	"github.com/paladignus/actajus/internal/module/person/application/command"
	"github.com/paladignus/actajus/internal/module/person/application/mapper"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

func TestPersonInputToDomain(t *testing.T) {
	t.Parallel()

	m := mapper.NewPersonMapper()
	got, err := m.PersonInputToDomain(command.CreatePersonCommand{
		Name:     "Jane Mary Doe",
		Birthday: "01/01/1990",
		Gender:   2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FirstName().Value() != "Jane" || got.LastName().Value() != "Mary Doe" {
		t.Fatalf("unexpected name mapping: first=%s last=%s", got.FirstName().Value(), got.LastName().Value())
	}
}

func TestPersonInputToDomainRequiresFullName(t *testing.T) {
	t.Parallel()

	m := mapper.NewPersonMapper()
	_, err := m.PersonInputToDomain(command.CreatePersonCommand{
		Name:     "Jane",
		Birthday: "01/01/1990",
		Gender:   2,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve sharedDomain.ValidationError
	if ok := AsValidationError(err, &ve); !ok {
		t.Fatalf("expected validation error, got %T", err)
	}
	if len(ve.Violations) != 1 || ve.Violations[0].Path != "name" || ve.Violations[0].Code != sharedDomain.CodeInvalid {
		t.Fatalf("unexpected violations: %+v", ve.Violations)
	}
}

func TestPersonInputToDomainRequiresName(t *testing.T) {
	t.Parallel()

	m := mapper.NewPersonMapper()
	_, err := m.PersonInputToDomain(command.CreatePersonCommand{
		Name:     "   ",
		Birthday: "01/01/1990",
		Gender:   2,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve sharedDomain.ValidationError
	if ok := AsValidationError(err, &ve); !ok {
		t.Fatalf("expected validation error, got %T", err)
	}
	if len(ve.Violations) != 1 || ve.Violations[0].Path != "name" || ve.Violations[0].Code != sharedDomain.CodeRequired {
		t.Fatalf("unexpected violations: %+v", ve.Violations)
	}
}

func AsValidationError(err error, target *sharedDomain.ValidationError) bool {
	ve, ok := err.(sharedDomain.ValidationError)
	if !ok {
		return false
	}
	*target = ve
	return true
}
