// Package mapper
package mapper

import (
	"time"

	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyMapper struct {
	addrMapper *addrMapper.AddressMapper
}

func NewCompanyMapper(addrMapper *addrMapper.AddressMapper) *CompanyMapper {
	return &CompanyMapper{
		addrMapper,
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

func (m *CompanyMapper) UpdateInputDomain(existing *domain.Company, input dto.UpdateCompanyRequest) error {
	return existing.UpdateBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		Apply()
}

func (m *CompanyMapper) DomainToOutput(projection CompanyProjection) dto.CompanyResponse {
	response := dto.CompanyResponse{
		ID:           projection.Company.ID(),
		Name:         projection.Company.Name().Value(),
		TradeName:    projection.Company.TradeName().Value(),
		CNPJ:         projection.Company.CNPJ().Value(),
		RegisteredBy: projection.Company.RegisteredBy(),
		CreatedAt:    projection.Company.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    projection.Company.UpdatedAt().Format(time.RFC3339),
	}

	if projection.Address != nil {
		addrResp := m.addrMapper.DomainToOutput(projection.Address)
		response.Address = &addrResp
	}
	return response
}

// func (m *CompanyMapper) DomainToOutput(company *domain.Company) dto.CompanyResponse {
// 	return dto.CompanyResponse{
// 		ID:           company.ID(),
// 		Name:         company.Name().Value(),
// 		TradeName:    company.TradeName().Value(),
// 		CNPJ:         company.CNPJ().Value(),
// 		RegisteredBy: company.RegisteredBy(),
// 		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
// 		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
// 	}
// }

// func (m *CompanyMapper) DomainListToOutput(companies []*domain.Company, page, pageSize int, total int64) dto.CompanyListResponse {
// 	responses := make([]dto.CompanyResponse, len(companies))
// 	for i, company := range companies {
// 		responses[i] = m.DomainToOutput(company)
// 	}
// 	return dto.CompanyListResponse{
// 		Companies: responses,
// 		Page:      page,
// 		PageSize:  pageSize,
// 		Total:     total,
// 	}
// }
