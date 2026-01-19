// Package mapper
package mapper

import (
	"time"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyMapper struct {
	// addressMapper
}

func NewCompanyMapper() *CompanyMapper {
	return &CompanyMapper{
		// addressMapper: NewAddressMapper()
	}
}

func (m *CompanyMapper) InputToDomain(input dto.CreateCompanyRequest) (*domain.Company, error) {
	// address, err := m.addressMappper.InputToDomain(input.Address)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid address", err)
	// }
	return domain.NewCompanyBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithRegisteredBy(input.RegisteredBy).
		Build()
	// if err != nil {
	// 	return nil, err
	// }
	// return company, nil
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
	return dto.CompanyResponse{
		ID:           company.ID(),
		Name:         company.Name().Value(),
		TradeName:    company.TradeName().Value(),
		CNPJ:         company.CNPJ().Value(),
		RegisteredBy: company.RegisteredBy(),
		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
	}
	// if company.CurrentAddress() != nil {
	// 	addr := m.addressMapper.DomainToOutput(company.CurrentAddress())
	// }
	// return response
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
