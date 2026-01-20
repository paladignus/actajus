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
	uow    domain.CompanyUnitOfWork
	mapper mapper.CompanyMapper
}

func NewCreateCompany(uow domain.CompanyUnitOfWork, mapper mapper.CompanyMapper) CreateCompany {
	return CreateCompany{uow, mapper}
}

func (c CreateCompany) Execute(ctx context.Context, input dto.CreateCompanyRequest) (*dto.CompanyResponse, error) {
	company, err := c.mapper.InputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	existing, err := c.uow.Company().FindByCNPJ(ctx, company.CNPJ().Value())
	if err != nil {
		return nil, fmt.Errorf("failed to check existing company: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrCNPJAlreadyExists
	}
	if err := c.uow.Begin(ctx); err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			c.uow.Rollback(ctx)
			panic(r)
		}
	}()
	if err := c.uow.Company().Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	if err := c.uow.Address().Create(ctx, company.Address()); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	if err := c.uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	response := c.mapper.DomainToOutput(company)
	return &response, nil
}
