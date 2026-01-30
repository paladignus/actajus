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
	uow        domain.CompanyUnitOfWork
	mapper     mapper.CompanyMapper
	projection mapper.CompanyProjectionMapper
}

func NewCreateCompany(
	uow domain.CompanyUnitOfWork,
	mapper mapper.CompanyMapper,
	projection mapper.CompanyProjectionMapper,
) CreateCompany {
	return CreateCompany{uow, mapper, projection}
}

func (c CreateCompany) Execute(
	ctx context.Context,
	input dto.CreateCompanyRequest,
) (*dto.CompanyReadModel, error) {
	company, err := c.mapper.CompanyInputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid company data: %w", err)
	}
	address, err := c.mapper.AddressInputToDomain(input.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address data: %w", err)
	}
	phone, err := c.mapper.PhoneInputToDomain(input.Phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone data: %w", err)
	}
	email, err := c.mapper.EmailInputToDomain(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email data: %w", err)
	}
	socialMedia, err := c.mapper.SocialMediaInputToDomain(input.SocialMedia)
	if err != nil {
		return nil, fmt.Errorf("invalid social media data: %w", err)
	}
	if err := c.uow.Begin(ctx); err != nil {
		return nil, err
	}
	defer c.uow.Rollback(ctx)
	if err := c.uow.Company().Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	if err := c.uow.Address().Create(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}
	if err := c.uow.Phone().Create(ctx, phone); err != nil {
		return nil, fmt.Errorf("failed to create phone: %w", err)
	}
	if err := c.uow.Email().Create(ctx, email); err != nil {
		return nil, fmt.Errorf("failed to create email: %w", err)
	}
	if err := c.uow.CompanyAddress().Create(ctx, company.ID(), address.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}
	if err := c.uow.CompanyPhone().Create(ctx, company.ID(), phone.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}
	if err := c.uow.CompanyEmail().Create(ctx, company.ID(), email.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}
	for _, sm := range socialMedia {
		if err := sm.SetCompanyID(company.ID()); err != nil {
			return nil, fmt.Errorf("failed to set company id: %w", err)
		}
		if err := c.uow.SocialMedia().Create(ctx, sm); err != nil {
			return nil, fmt.Errorf("failed to create social media: %w", err)
		}
	}
	if err := c.uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}
	return c.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		email,
		socialMedia,
	), nil
}
