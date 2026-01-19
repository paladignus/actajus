// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CreateCompany struct {
	repository domain.CompanyRepository
	mapper     mapper.CompanyMapper
}

func NewCreateCompany(repository domain.CompanyRepository, mapper mapper.CompanyMapper) CreateCompany {
	return CreateCompany{repository, mapper}
}

func (c CreateCompany) Execute(ctx context.Context, input dto.CreateCompanyRequest) (*dto.CompanyResponse, error) {
	company, err := c.mapper.InputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	existing, err := c.repository.FindByCNPJ(ctx, company.CNPJ().Value())
	if err != nil {
		return nil, fmt.Errorf("failed to check existing company: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrCNPJAlreadyExists
	}
	if err := c.repository.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	response := c.mapper.DomainToOutput(company)
	return &response, nil
}
