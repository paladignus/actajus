// Package usecase
package usecase

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
)

type FindByID struct {
	repository repository.Factory
	projection mapper.CompanyProjectionMapper
}

func NewFindByID(
	repository repository.Factory,
	projection mapper.CompanyProjectionMapper,
) FindByID {
	return FindByID{
		repository,
		projection,
	}
}

func (u FindByID) Execute(ctx context.Context, id int64) (*readmodel.CompanyReadModel, error) {
	company, err := u.repository.Company().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	address, err := u.repository.Address().FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	phone, err := u.repository.Phone().FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	email, err := u.repository.Email().FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	socialMedia, err := u.repository.SocialMedia().FindByIDCompany(ctx, id)
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
