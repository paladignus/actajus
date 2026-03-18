// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
)

type FindByCNPJ struct {
	repository repository.Factory
	projection mapper.CompanyProjectionMapper
}

func NewFindByCNPJ(
	repository repository.Factory,
	projection mapper.CompanyProjectionMapper,
) FindByCNPJ {
	return FindByCNPJ{
		repository,
		projection,
	}
}

func (u FindByCNPJ) Execute(ctx context.Context, cnpj string) (*dto.CompanyReadModel, error) {
	company, err := u.repository.Company().FindByCNPJ(ctx, cnpj)
	if err != nil {
		return nil, err
	}
	address, err := u.repository.Address().FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	phone, err := u.repository.Phone().FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	email, err := u.repository.Email().FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	socialMedia, err := u.repository.SocialMedia().FindByIDCompany(ctx, company.ID())
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
