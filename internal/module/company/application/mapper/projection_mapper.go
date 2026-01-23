// Package mapper
package mapper

import (
	"time"

	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailMapper "github.com/paladignus/actajus/internal/module/email/application/mapper"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaMapper "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type CompanyProjectionMapper struct {
	addrMapper        addrMapper.AddressPrejectionMapper
	phoneMapper       phoneMapper.PhoneProjectionMapper
	emailMapper       emailMapper.EmailProjectionMapper
	socialMediaMapper socialMediaMapper.SocialMediaProjectionMapper
}

func NewCompanyProjectionMapper(
	addrMapper addrMapper.AddressPrejectionMapper,
	phoneMapper phoneMapper.PhoneProjectionMapper,
	emailMapper emailMapper.EmailProjectionMapper,
	socialMediaMapper socialMediaMapper.SocialMediaProjectionMapper,
) CompanyProjectionMapper {
	return CompanyProjectionMapper{
		addrMapper,
		phoneMapper,
		emailMapper,
		socialMediaMapper,
	}
}

func (m *CompanyProjectionMapper) ProjectCompanyToReadModel(
	company *domain.Company,
	address *addrDomain.Address,
	phone *phoneDomain.Phone,
	email *emailDomain.Email,
	socialMedia []*socialMediaDomain.SocialMedia,
) *dto.CompanyReadModel {
	readModel := &dto.CompanyReadModel{
		ID:           company.ID(),
		Name:         company.Name().Value(),
		TradeName:    company.TradeName().Value(),
		CNPJ:         company.CNPJ().Value(),
		RegisteredBy: company.RegisteredBy(),
		CreatedAt:    company.CreatedAt().Format(time.RFC3339),
		UpdatedAt:    company.UpdatedAt().Format(time.RFC3339),
	}
	if address != nil {
		readModel.Address = m.addrMapper.ProjectAddressToReadModel(address)
	}
	if phone != nil {
		readModel.Phone = m.phoneMapper.ProjectPhoneToReadModel(phone)
	}
	if email != nil {
		readModel.Email = m.emailMapper.ProjectEmailToReadModel(email)
	}
	if len(socialMedia) > 0 {
		for _, socialMedia := range socialMedia {
			readModel.SocialMedia = append(readModel.SocialMedia, m.socialMediaMapper.ProjectSocialMediaToReadModel(socialMedia))
		}
	}
	return readModel
}
