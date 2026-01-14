// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type Enterprise struct {
	uow repository.UnitOfWorkEnterprise
}

func NewEnterprise(
	uow repository.UnitOfWorkEnterprise,
) Enterprise {
	return Enterprise{uow}
}

func (e Enterprise) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	if err := e.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer e.uow.Rollback(ctx)
	enterprise, err := entity.NewEnterprise(input.Enterprise)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	idEnterprise, err := e.uow.Enterprise().Create(ctx, enterprise)
	if err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}
	address, err := entity.NewAddress(input.Address)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}
	idAddress, err := e.uow.Address().Create(ctx, address)
	if err != nil {
		return fmt.Errorf("database error while saving address: %w", err)
	}
	if err = e.uow.AddressEnterprise().Create(ctx, idAddress, idEnterprise); err != nil {
		return fmt.Errorf("database error while saving address enterprise relation: %w", err)
	}
	return e.uow.Commit(ctx)
}
