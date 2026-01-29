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
	if err := u.uow.Begin(ctx); err != nil {
		return nil, err
	}
	defer u.uow.Rollback(ctx)
	company, err := u.uow.Company().FindByID(ctx, input.IDCompany)
	if err != nil {
		return nil, err
	}
	err = u.mapper.UpdateInputDomain(company, input)
	if err != nil {
		return nil, err
	}
	if err := u.uow.Company().Update(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
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
