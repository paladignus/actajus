// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type UpdateCompany struct {
	uow        domain.CompanyUnitOfWork
	mapper     mapper.CompanyMapper
	projection mapper.CompanyProjectionMapper
}

func NewUpdateCompany(
	uow domain.CompanyUnitOfWork,
	mapper mapper.CompanyMapper,
	projection mapper.CompanyProjectionMapper,
) UpdateCompany {
	return UpdateCompany{uow, mapper, projection}
}

func (u UpdateCompany) Execute(ctx context.Context, input dto.UpdateCompanyRequest) (*dto.CompanyReadModel, error) {
	company, err := u.mapper.UpdateInputDomain(input)
	if err != nil {
		return nil, err
	}
	address, err := u.mapper.UpdateAddressInputToDomain(input.Address)
	if err != nil {
		return nil, err
	}
	phone, err := u.mapper.UpdatePhoneInputToDomain(input.Phone)
	if err != nil {
		return nil, err
	}
	email, err := u.mapper.UpdateEmailInputToDomain(input.Email)
	if err != nil {
		return nil, err
	}
	socialMedia, err := u.mapper.UpdateSocialMediaInputToDomain(input.SocialMedia)
	if err != nil {
		return nil, fmt.Errorf("invalid social media data to company: %w", err)
	}
	if err := u.uow.Begin(ctx); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}
	defer u.uow.Rollback(ctx)

	// Delete old relationships before updating
	if err := u.uow.CompanyAddress().DeleteByIDCompany(ctx, company.ID()); err != nil {
		return nil, fmt.Errorf("failed to delete company address relationship: %w", err)
	}
	if err := u.uow.CompanyPhone().DeleteByIDCompany(ctx, company.ID()); err != nil {
		return nil, fmt.Errorf("failed to delete company phone relationship: %w", err)
	}
	if err := u.uow.CompanyEmail().DeleteByIDCompany(ctx, company.ID()); err != nil {
		return nil, fmt.Errorf("failed to delete company email relationship: %w", err)
	}

	if err := u.uow.Company().Update(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}
	if err := u.uow.Address().Update(ctx, *address); err != nil {
		return nil, fmt.Errorf("failed to update address to company: %w", err)
	}
	if err := u.uow.Phone().Update(ctx, *phone); err != nil {
		return nil, fmt.Errorf("failed to upate phone to company %w", err)
	}
	if err := u.uow.Email().Update(ctx, *email); err != nil {
		return nil, fmt.Errorf("failed to update email to company: %w", err)
	}

	// Create new relationships
	if err := u.uow.CompanyAddress().Create(ctx, company.ID(), address.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship company address: %w", err)
	}
	if err := u.uow.CompanyPhone().Create(ctx, company.ID(), phone.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship company phone: %w", err)
	}
	if err := u.uow.CompanyEmail().Create(ctx, company.ID(), email.ID()); err != nil {
		return nil, fmt.Errorf("failed to create relationship company email: %w", err)
	}

	for _, sm := range socialMedia {
		if err := u.uow.SocialMedia().Update(ctx, *sm); err != nil {
			return nil, fmt.Errorf("failed to update social media to company: %w", err)
		}
	}
	if err := u.uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}
	return u.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		email,
		socialMedia,
	), nil
}
