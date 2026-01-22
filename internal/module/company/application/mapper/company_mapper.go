// Package mapper
package mapper

import (
	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	emailMapper "github.com/paladignus/actajus/internal/module/email/application/mapper"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
)

type CompanyMapper struct {
	addrMapper  *addrMapper.AddressMapper
	phoneMapper *phoneMapper.PhoneMapper
	emailMapper *emailMapper.EmailMapper
}

func NewCompanyMapper(
	addrMapper *addrMapper.AddressMapper,
	phoneMapper *phoneMapper.PhoneMapper,
	emailMapper *emailMapper.EmailMapper,
) *CompanyMapper {
	return &CompanyMapper{
		addrMapper,
		phoneMapper,
		emailMapper,
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

func (m *CompanyMapper) EmailInputToDomain(input emailDTO.CreateEmailRequest) (*emailDomain.Email, error) {
	return m.emailMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateInputDomain(existing *domain.Company, input dto.UpdateCompanyRequest) error {
	return existing.UpdateBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		Apply()
}
