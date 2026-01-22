// Package mapper
package mapper

import (
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
)

type CompanyMapper struct {
	addrMapper  *addrMapper.AddressMapper
	phoneMapper *phoneMapper.PhoneMapper
}

func NewCompanyMapper(
	addrMapper *addrMapper.AddressMapper,
	phphoneMapper *phoneMapper.PhoneMapper,
) *CompanyMapper {
	return &CompanyMapper{
		addrMapper,
		phphoneMapper,
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
	return m.addrMapper.InputToDomain(input)
}

func (m *CompanyMapper) PhoneInputToDomain(input phoneDTO.CreatePhoneRequest) (*phoneDomain.Phone, error) {
	return m.phoneMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateInputDomain(existing *domain.Company, input dto.UpdateCompanyRequest) error {
	return existing.UpdateBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		Apply()
}
