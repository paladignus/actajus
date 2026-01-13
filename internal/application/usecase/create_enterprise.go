// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
)

type Enterprise struct {
	uow postgres.UnitOfWorkEnterprise
}

func NewEnterprise(
	uow postgres.UnitOfWorkEnterprise,
) Enterprise {
	return Enterprise{uow}
}

func (e Enterprise) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	if err := e.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer e.uow.Rollback(ctx)

	// Get repositories with transaction-aware pool
	// persistence := postgres.NewPersistence(e.uow.GetPgxPool())
	// enterpriseRepo := persistence.Enterprise()
	// addressRepo := persistence.Address()

	enterprise, err := entity.NewEnterprise(input.Enterprise)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	_, err = e.uow.Enterprise().Create(ctx, enterprise)
	if err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}

	address, err := entity.NewAddress(input.Address)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	_, err = e.uow.Address().Create(ctx, address)
	if err != nil {
		return fmt.Errorf("database error while saving address: %w", err)
	}
	return e.uow.Commit(ctx)
}
