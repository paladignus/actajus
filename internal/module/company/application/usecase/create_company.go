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
	company, err := c.mapper.CompanyInputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid company data: %w", err)
	}
	address, err := c.mapper.AddressInputToDomain(input.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address data: %w", err)
	}
	existing, err := c.uow.Company().FindByCNPJ(ctx, company.CNPJ().Value())
	if err != nil {
		return nil, err
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
	if err := c.uow.Address().Create(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	addressID := address.ID()
	company.SetAddress(addressID)
	if err := c.uow.Company().Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	if err := c.uow.CompanyAddress().Create(ctx, company.ID(), addressID); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}
	if err := c.uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}
	response := c.mapper.DomainToOutput(
		mapper.CompanyProjection{Company: company, Address: address},
	)
	return &response, nil
}
