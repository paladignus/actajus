// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type EnterpriseUpdate struct {
	uow repository.UnitOfWorkEnterprise
}

func NewEnterpriseUpdate(uow repository.UnitOfWorkEnterprise) EnterpriseUpdate {
	return EnterpriseUpdate{uow}
}

func (e EnterpriseUpdate) Execute(ctx context.Context, input dto.EnterpriseInput) error {
	return nil
}
