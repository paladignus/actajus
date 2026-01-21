// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyProjectionMapper struct {
	mapper mapper.AddressPrejectionMapper
}

func NewCompanyProjectionMapper(mapper mapper.AddressPrejectionMapper) CompanyProjectionMapper {
	return CompanyProjectionMapper{mapper}
}

func (m *CompanyProjectionMapper) ProjectCompanyToReadModel(
	company *domain.Company,
	address *addrDomain.Address,
) *dto.CompanyReadModel {
	readModel := &dto.CompanyReadModel{
		ID:           company.ID(),
		Name:         company.Name().Value(),
		TradeName:    company.TradeName().Value(),
		CNPJ:         company.CNPJ().Value(),
		RegisteredBy: company.RegisteredBy(),
		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
	}
	if address != nil {
		readModel.Address = m.mapper.ProjectAddressToReadModel(address)
	}
	return readModel
}
