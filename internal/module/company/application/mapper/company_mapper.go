// Package mapper
package mapper

import (
	"fmt"
	"time"

	addressMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyMapper struct {
	addressMapper *addressMapper.AddressMapper
}

func NewCompanyMapper(addressMapper *addressMapper.AddressMapper) *CompanyMapper {
	return &CompanyMapper{
		addressMapper,
	}
}

func (m *CompanyMapper) InputToDomain(input dto.CreateCompanyRequest) (*domain.Company, error) {
	address, err := m.addressMapper.InputToDomain(input.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address %w", err)
	}
	return domain.NewCompanyBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithRegisteredBy(input.RegisteredBy).
		WithAddress(address).
		Build()
}

func (m *CompanyMapper) UpdateInputDomain(existing *domain.Company, input dto.UpdateCompanyRequest) error {
	return existing.UpdateBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		Apply()
}

// func (m *CompanyMapper) CreateAddressInputToDomain(input dto.ChangeAddressRequest) (*addressDomain.Address, error) {
// 	return m.addressMapper.InputToDomain(input.Address)
// }

func (m *CompanyMapper) DomainToOutput(company *domain.Company) dto.CompanyResponse {
	response := dto.CompanyResponse{
		ID:           company.ID(),
		Name:         company.Name().Value(),
		TradeName:    company.TradeName().Value(),
		CNPJ:         company.CNPJ().Value(),
		RegisteredBy: company.RegisteredBy(),
		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
	}
	if company.Address() != nil {
		addr := m.addressMapper.DomainToOutput(company.Address())
		response.Address = &addr
	}
	return response
}

// func (m *CompanyMapper) DomainToOutputWithHistory(company *domain.Company) dto.CompanyWithHistoryResponse {
// 	response := dto.CompanyWithHistoryResponse{
// 		ID:           company.ID(),
// 		Name:         company.Name().Value(),
// 		TradeName:    company.TradeName().Value(),
// 		CNPJ:         company.CNPJ().Value(),
// 		RegisteredBy: company.RegisteredBy(),
// 		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
// 		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
// 	}
// if company.CurrentAddress() != nil {
// 	addr := m.addressMapper.DomainToOutput(company.CurrentAddress())
// 	response.CurrentAddress = &addr
// }
// if len(company.AddressHistory()) > 0 {
// 	response.AddressHistory = make([]addressDTO.AddressResponse, len(company.AddressHistory()))
// 	for i, addr := range company.AddressHistory() {
// 		response.AddressHistory[i] = m.addressMapper.DomainToOutput(&addr)
// 	}
// }
// 	return response
// }

func (m *CompanyMapper) DomainListToOutput(companies []*domain.Company, page, pageSize int, total int64) dto.CompanyListResponse {
	responses := make([]dto.CompanyResponse, len(companies))
	for i, company := range companies {
		responses[i] = m.DomainToOutput(company)
	}
	return dto.CompanyListResponse{
		Companies: responses,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
	}
}
