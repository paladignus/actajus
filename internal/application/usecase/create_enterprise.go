// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/domain/repository"
	unitofwork "github.com/paladignus/actajus/internal/domain/unit_of_work"
)

type Enterprise struct {
	uow         unitofwork.IUnitOfWork
	persistence repository.IEnterprise
}

func NewEnterprise(uow unitofwork.IUnitOfWork, persistence repository.IEnterprise) Enterprise {
	return Enterprise{uow, persistence}
}

func (e Enterprise) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	enterprise, err := entity.NewEnterprise(input)
	if err != nil {
		return fmt.Errorf("use case create enterprise, invalid input: %w", err)
	}

	if err := e.uow.Begin(ctx); err != nil {
		return fmt.Errorf("database error for begin transaction: %w", err)
	}
	defer e.uow.Rollback(ctx)

	id, err := e.persistence.Create(ctx, enterprise)
	fmt.Println(id)
	if err != nil {
		return fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return e.uow.Commit(ctx)
}
