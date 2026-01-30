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
	// email, err := u.mapper.EmailInputToDomain(input.Email)
	// if err != nil {
	//   return nil, err
	// }
	// socialMedia, err := u.mapper.SocialMediaInputToDomain(input.SocialMedia)
	// if err != nil {
	//   return nil, err
	// }
	if err := u.uow.Begin(ctx); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}
	defer u.uow.Rollback(ctx)
	if err := u.uow.Company().Update(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}
	if err := u.uow.Address().Update(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address to company: %w", err)
	}
	if err := u.uow.Phone().Update(ctx, phone); err != nil {
		return nil, fmt.Errorf("failed to upate phone to company %w", err)
	}
	if err := u.uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}
	return nil, nil
	// return u.projection.ProjectCompanyToReadModel(
	// 	company,
	// 	address,
	// 	phone,
	// 	email,
	// 	socialMedia,
	// ), nil
}
