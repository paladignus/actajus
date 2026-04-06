// Package mapper
package mapper

import (
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/readmodel"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type CompanyProjectionMapper struct{}

func NewCompanyProjectionMapper() CompanyProjectionMapper {
	return CompanyProjectionMapper{}
}

func (m *CompanyProjectionMapper) ProjectCompanyToReadModel(
	company *domain.Company,
	address *addrDomain.Address,
	phone *phoneDomain.Phone,
	email *emailDomain.Email,
	socialMedia []*socialMediaDomain.SocialMedia,
) *readmodel.CompanyReadModel {
	readModel := &readmodel.CompanyReadModel{
		ID:               company.ID().Value(),
		Name:             company.Name().Value(),
		TradeName:        company.TradeName().Value(),
		CNPJ:             company.CNPJ().Value(),
		RegisteredByName: company.RegisteredByName().Value(),
		CreatedAt:        company.CreatedAt(),
		UpdatedAt:        company.UpdatedAt(),
	}
	if address != nil {
		readModel.Address = &readmodel.CompanyAddressReadModel{
			ID:           address.ID().Value(),
			ZIP:          address.ZIP().Value(),
			Title:        address.Title().Value(),
			Street:       address.Street().Value(),
			Number:       address.Number(),
			Neighborhood: address.Neighborhood().Value(),
			City:         address.City().Value(),
			State:        address.State().Value(),
			Country:      address.Country().Value(),
			CreatedAt:    address.CreatedAt(),
			UpdatedAt:    address.UpdatedAt(),
		}
		if value := address.Complement().Value(); value != "" {
			readModel.Address.Complement = &value
		}
		if value := address.Reference().Value(); value != "" {
			readModel.Address.Reference = &value
		}
	}
	if phone != nil {
		readModel.Phones = append(readModel.Phones, &readmodel.CompanyPhoneReadModel{
			ID:         phone.ID().Value(),
			Number:     phone.Number().Value(),
			Kind:       phone.Kind().Value(),
			Department: phone.Department().Value(),
			CreatedAt:  phone.CreatedAt(),
			UpdatedAt:  phone.UpdatedAt(),
		})
	}
	if email != nil {
		readModel.Emails = append(readModel.Emails, &readmodel.CompanyEmailReadModel{
			ID:        email.ID().Value(),
			Address:   email.Address().Value(),
			CreatedAt: email.CreatedAt(),
			UpdatedAt: email.UpdatedAt(),
		})
	}
	if len(socialMedia) > 0 {
		for _, sm := range socialMedia {
			readModel.SocialMedia = append(readModel.SocialMedia, &readmodel.CompanySocialMediaReadModel{
				ID:        sm.ID().Value(),
				IDCompany: sm.IDCompany().Value(),
				Platform:  sm.Platform().Value(),
				URL:       sm.URL().Value(),
				CreatedAt: sm.CreatedAt(),
				UpdatedAt: sm.UpdatedAt(),
			})
		}
	}
	return readModel
}
