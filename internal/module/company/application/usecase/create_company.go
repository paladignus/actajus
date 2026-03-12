// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type CreateCompany struct {
	uow        uow.UnitOfWork
	repository repository.Factory
	mapper     mapper.CompanyMapper
	projection mapper.CompanyProjectionMapper
}

func NewCreateCompany(
	uow uow.UnitOfWork,
	repository repository.Factory,
	mapper mapper.CompanyMapper,
	projection mapper.CompanyProjectionMapper,
) CreateCompany {
	return CreateCompany{uow, repository, mapper, projection}
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
	err = c.uow.Do(ctx, func(tx uow.Tx) error {
		r := c.repository.WithTx(tx)
		if err := r.Company().Create(ctx, company); err != nil {
			return fmt.Errorf("failed to create company: %w", err)
		}
		if err := r.Address().Create(ctx, address); err != nil {
			return fmt.Errorf("failed to create address: %w", err)
		}
		if err := r.Phone().Create(ctx, phone); err != nil {
			return fmt.Errorf("failed to create phone: %w", err)
		}
		if err := r.Email().Create(ctx, email); err != nil {
			return fmt.Errorf("failed to create email: %w", err)
		}
		if err := r.CompanyAddress().Create(ctx, company.ID(), address.ID()); err != nil {
			return fmt.Errorf("failed to create relationship company address: %w", err)
		}
		if err := r.CompanyPhone().Create(ctx, company.ID(), phone.ID()); err != nil {
			return fmt.Errorf("failed to create relationship company phone: %w", err)
		}
		if err := r.CompanyEmail().Create(ctx, company.ID(), email.ID()); err != nil {
			return fmt.Errorf("failed to create relationship company email: %w", err)
		}
		for _, sm := range socialMedia {
			if err := sm.SetCompanyID(company.ID()); err != nil {
				return fmt.Errorf("failed to set company id: %w", err)
			}
			if err := r.SocialMedia().Create(ctx, sm); err != nil {
				return fmt.Errorf("failed to create social media: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return c.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		email,
		socialMedia,
	), nil
}
