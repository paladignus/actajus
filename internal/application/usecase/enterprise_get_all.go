// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type EnterpriseGetAll struct {
	uow repository.UnitOfWorkEnterprise
}

func NewEnterpriseGetAll(uow repository.UnitOfWorkEnterprise) EnterpriseGetAll {
	return EnterpriseGetAll{uow}
}

func (e EnterpriseGetAll) Execute(ctx context.Context) (enterprises []dto.EnterpriseOutput, err error) {
	data, err := e.uow.Enterprise().GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error while getting all enterprises: %w", err)
	}
	for _, enterprise := range data {
		enterprises = append(enterprises, dto.EnterpriseOutput{
			IDEnterprise: enterprise.IDEnterprise,
			RegisteredBy: enterprise.RegisteredBy,
			Name:         enterprise.Name.Value(),
			TradeName:    enterprise.TradeName.Value(),
			CNPJ:         enterprise.CNPJ.Value(),
		})
	}
	return enterprises, nil
}
