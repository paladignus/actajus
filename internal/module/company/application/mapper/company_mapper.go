// Package mapper
package mapper

import (
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	"github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyMapper struct {
	mapper *mapper.AddressMapper
}

func NewCompanyMapper(
	mapper *mapper.AddressMapper,
) *CompanyMapper {
	return &CompanyMapper{
		mapper,
	}
}

func (m *CompanyMapper) CompanyInputToDomain(input dto.CreateCompanyRequest) (*domain.Company, error) {
	return domain.NewCompanyBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithRegisteredBy(input.RegisteredBy).
		Build()
}

func (m *CompanyMapper) AddressInputToDomain(input addrDTO.CreateAddressRequest) (*addrDomain.Address, error) {
	return m.mapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateInputDomain(existing *domain.Company, input dto.UpdateCompanyRequest) error {
	return existing.UpdateBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		Apply()
}
