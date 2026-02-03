// Package usecase
package usecase

import (
	"context"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/domain"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
)

type FindByID struct {
	company    domain.CompanyRepository
	address    addrDomain.AddressRepository
	phone      phoneDomain.PhoneRepository
	projection mapper.CompanyProjectionMapper
}

func NewFindByID(
	company domain.CompanyRepository,
	address addrDomain.AddressRepository,
	phone phoneDomain.PhoneRepository,
	projection mapper.CompanyProjectionMapper,
) FindByID {
	return FindByID{
		company,
		address,
		phone,
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
	phone, err := f.phone.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	return f.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		nil, nil), nil
}
