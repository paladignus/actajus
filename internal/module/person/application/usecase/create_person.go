// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/person/application/command"
	"github.com/paladignus/actajus/internal/module/person/application/mapper"
	"github.com/paladignus/actajus/internal/module/person/application/readmodel"
	"github.com/paladignus/actajus/internal/module/person/domain"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/application/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
)

type CreatePerson struct {
	person     domain.PersonRepository
	mapper     mapper.PersonMapper
	projection mapper.PersonProjection
	validator  *validation.Validator
}

func NewCreatePerson(
	person domain.PersonRepository,
	mapper mapper.PersonMapper,
	projection mapper.PersonProjection,
) CreatePerson {
	return CreatePerson{
		person:     person,
		mapper:     mapper,
		projection: projection,
		validator:  validation.New(),
	}
}

func (c CreatePerson) Execute(ctx context.Context, input command.CreatePersonCommand) (*readmodel.PersonReadModel, error) {
	// Validação sintática do DTO
	vs := c.validator.ValidateStruct(input)
	if err := sharedAdapter.ViolationsToDomainError(vs); err != nil {
		return nil, err
	}

	// Mapeamento e validação de domínio
	person, err := c.mapper.PersonInputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("error mapping input to domain: %w", err)
	}
	if err := c.person.Create(ctx, person); err != nil {
		return nil, fmt.Errorf("error creating person: %w", err)
	}
	return c.projection.ProjectPersonReadModel(person), nil
}
