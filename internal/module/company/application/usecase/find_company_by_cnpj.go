// Package usecase
package usecase

import (
	"context"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/application/mapper"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type FindByCNPJ struct {
	company     domain.CompanyRepository
	address     addrDomain.AddressRepository
	phone       phoneDomain.PhoneRepository
	email       emailDomain.EmailRepository
	socialMedia socialMediaDomain.SocialMediaRepository
	projection  mapper.CompanyProjectionMapper
}

func NewFindByCNPJ(
	company domain.CompanyRepository,
	address addrDomain.AddressRepository,
	phone phoneDomain.PhoneRepository,
	email emailDomain.EmailRepository,
	socialMedia socialMediaDomain.SocialMediaRepository,
	projection mapper.CompanyProjectionMapper,
) FindByCNPJ {
	return FindByCNPJ{
		company,
		address,
		phone,
		email,
		socialMedia,
		projection,
	}
}

func (f FindByCNPJ) Execute(ctx context.Context, cnpj string) (*dto.CompanyReadModel, error) {
	company, err := f.company.FindByCNPJ(ctx, cnpj)
	if err != nil {
		return nil, err
	}
	address, err := f.address.FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	phone, err := f.phone.FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	email, err := f.email.FindByIDCompany(ctx, company.ID())
	if err != nil {
		return nil, err
	}
	socialMedia, err := f.socialMedia.FindByIDCompany(ctx, company.ID())
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
