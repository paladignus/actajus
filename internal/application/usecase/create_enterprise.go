// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type Enterprise struct {
	uow         postgres.IUnitOfWork
	persistence repository.IEnterprise
}

func NewEnterprise(
	uow postgres.IUnitOfWork,
	persistence repository.IEnterprise,
) Enterprise {
	return Enterprise{uow, persistence}
}

func (e Enterprise) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	if err := e.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer e.uow.Rollback(ctx)
	enterprise, err := entity.NewEnterprise(input)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	id, err := e.persistence.Create(ctx, enterprise)
	fmt.Println(id)
	if err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return e.uow.Commit(ctx)
}
