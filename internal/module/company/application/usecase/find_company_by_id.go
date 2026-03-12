// Package usecase
package usecase

import (
	"context"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type FindByID struct {
	company     repository.CompanyRepository
	address     addrDomain.AddressRepository
	phone       phoneDomain.PhoneRepository
	email       emailDomain.EmailRepository
	socialMedia socialMediaDomain.SocialMediaRepository
	projection  mapper.CompanyProjectionMapper
}

func NewFindByID(
	company repository.CompanyRepository,
	address addrDomain.AddressRepository,
	phone phoneDomain.PhoneRepository,
	email emailDomain.EmailRepository,
	socialMedia socialMediaDomain.SocialMediaRepository,
	projection mapper.CompanyProjectionMapper,
) FindByID {
	return FindByID{
		company,
		address,
		phone,
		email,
		socialMedia,
		projection,
	}
}

func (f FindByID) Execute(ctx context.Context, id uint) (*dto.CompanyReadModel, error) {
	company, err := f.company.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	address, err := f.address.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	phone, err := f.phone.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	email, err := f.email.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	socialMedia, err := f.socialMedia.FindByIDCompany(ctx, id)
	if err != nil {
		return nil, err
	}
	return f.projection.ProjectCompanyToReadModel(
		company,
		address,
		phone,
		email,
		socialMedia,
	), nil
}
