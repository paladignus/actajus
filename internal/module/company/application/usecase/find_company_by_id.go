// Package usecase
package usecase

import (
	"context"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type FindByID struct {
	company    domain.CompanyRepository
	address    addrDomain.AddressRepository
	projection mapper.CompanyProjectionMapper
}

func NewFindByID(
	company domain.CompanyRepository,
	address addrDomain.AddressRepository,
	projection mapper.CompanyProjectionMapper,
) FindByID {
	return FindByID{
		company,
		address,
		projection,
	}
}

func (f FindByID) Execute(ctx context.Context, id uint) (*dto.CompanyReadModel, error) {
	company, err := f.company.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	address, err := f.address.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	return f.projection.ProjectCompanyToReadModel(
		company,
		address,
		nil, nil, nil), nil
}
