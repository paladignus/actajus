// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type DeleteCompany struct {
	uow repository.UnitOfWorkCompany
}

func NewDeleteCompany(uow repository.UnitOfWorkCompany) DeleteCompany {
	return DeleteCompany{uow}
}

func (d DeleteCompany) Execute(ctx context.Context, idCompany uint) error {
	if err := d.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer d.uow.Rollback(ctx)
	if err := d.uow.Company().Delete(ctx, idCompany); err != nil {
		return fmt.Errorf("database error while delete company: %w", err)
	}
	return d.uow.Commit(ctx)
}
