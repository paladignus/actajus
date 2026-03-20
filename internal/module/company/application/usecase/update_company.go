// Package usecase
package usecase

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/command"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type UpdateCompany struct {
	uow        uow.UnitOfWork
	repository repository.Factory
	mapper     mapper.CompanyMapper
	projection mapper.CompanyProjectionMapper
}

func NewUpdateCompany(
	uow uow.UnitOfWork,
	repository repository.Factory,
	mapper mapper.CompanyMapper,
	projection mapper.CompanyProjectionMapper,
) UpdateCompany {
	return UpdateCompany{uow, repository, mapper, projection}
}

func (u UpdateCompany) Execute(
	ctx context.Context,
	input command.UpdateCompanyCommand,
) (*readmodel.CompanyReadModel, error) {
	company, err := u.mapper.UpdateInputToDomain(input)
	if err != nil {
		return nil, fmt.Errorf("invalid company data: %w", err)
	}
	address, err := u.mapper.UpdateAddressInputToDomain(input.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address data: %w", err)
	}
	phone, err := u.mapper.UpdatePhoneInputToDomain(input.Phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone data: %w", err)
	}
	email, err := u.mapper.UpdateEmailInputToDomain(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email data: %w", err)
	}
	socialMedia, err := u.mapper.UpdateSocialMediaInputToDomain(input.SocialMedia)
	if err != nil {
		return nil, fmt.Errorf("invalid social media data: %w", err)
	}
	err = u.uow.Do(ctx, func(tx uow.Tx) error {
		r := u.repository.WithTx(tx)

		if err := r.Company().Update(ctx, company); err != nil {
			return fmt.Errorf("failed to update company: %w", err)
		}
		if err := r.Address().Update(ctx, *address); err != nil {
			return fmt.Errorf("failed to update address: %w", err)
		}
		if err := r.Phone().Update(ctx, *phone); err != nil {
			return fmt.Errorf("failed to update phone: %w", err)
		}
		if err := r.Email().Update(ctx, *email); err != nil {
			return fmt.Errorf("failed to update email: %w", err)
		}
		for _, sm := range socialMedia {
			if err := r.SocialMedia().Update(ctx, *sm); err != nil {
				return fmt.Errorf("failed to update social media to company: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		email,
		socialMedia,
	), nil
}
