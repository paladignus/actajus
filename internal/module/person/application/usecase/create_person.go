// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/person/application/dto"
	"github.com/paladignus/actajus/internal/module/person/application/mapper"
	"github.com/paladignus/actajus/internal/module/person/domain"
)

type CreatePerson struct {
	person     domain.PersonRepository
	mapper     mapper.PersonMapper
	projection mapper.PersonProjection
}

func NewCreatePerson(
	person domain.PersonRepository,
	mapper mapper.PersonMapper,
	projection mapper.PersonProjection,
) CreatePerson {
	return CreatePerson{
		person,
		mapper,
		projection,
	}
}

func (c CreatePerson) Execute(ctx context.Context, input dto.CreatePersonRequest) (*dto.PersonReadModel, error) {
	person, err := c.mapper.PersonInputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("error mapping input to domain: %w", err)
	}
	if err := c.person.Create(ctx, person); err != nil {
		return nil, fmt.Errorf("error creating person: %w", err)
	}
	return c.projection.ProjectPersonReadModel(person), nil
}
