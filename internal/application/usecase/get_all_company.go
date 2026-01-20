// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/application/dto"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type GetAllCompany struct {
	uow repository.UnitOfWorkCompany
}

func NewGetAllCompany(uow repository.UnitOfWorkCompany) GetAllCompany {
	return GetAllCompany{uow}
}

func (g GetAllCompany) Execute(ctx context.Context) ([]dto.CompanyInputOutput, error) {
	companies, err := g.uow.Company().GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error while getting all enterprises: %w", err)
	}
	return companies, nil
}
