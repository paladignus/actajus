// Package mapper
package mapper

import (
	"time"

	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/command"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type CompanyMapper struct{}

func NewCompanyMapper() *CompanyMapper {
	return &CompanyMapper{}
}

func (m *CompanyMapper) CompanyInputToDomain(input command.CreateCompanyCommand) (*domain.Company, error) {
	// v := validation.New(validation.PT)
	// if err := v.ValidateStruct(input); err != nil {
	// 	return nil, err
	// }
	return domain.NewCompanyBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithRegisteredBy(input.RegisteredBy).
		Build()
}

func (m *CompanyMapper) UpdateInputToDomain(input command.UpdateCompanyCommand) (*domain.Company, error) {
	// v := validation.New(validation.PT)
	// if err := v.ValidateStruct(input); err != nil {
	// 	return nil, err
	// }
	return domain.NewCompanyBuilder().
		WithID(input.IDCompany).
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithUpdatedAt(time.Now()).
		Build()
}

func (m *CompanyMapper) AddressInputToDomain(input command.CreateCompanyAddressCommand) (*addrDomain.Address, error) {
	return addrDomain.NewAddressBuilder().
		WithZIP(input.ZIP).
		WithTitle(input.Title).
		WithStreet(input.Street).
		WithNumber(input.Number).
		WithComplement(input.Complement).
		WithReference(input.Reference).
		WithNeighborhood(input.Neighborhood).
		WithCity(input.City).
		WithState(input.State).
		WithCountry(input.Country).
		Build()
}

func (m *CompanyMapper) UpdateAddressInputToDomain(input command.UpdateCompanyAddressCommand) (*addrDomain.Address, error) {
	return addrDomain.NewAddressBuilder().
		WithID(input.IDAddress).
		WithZIP(input.ZIP).
		WithTitle(input.Title).
		WithStreet(input.Street).
		WithNumber(input.Number).
		WithComplement(input.Complement).
		WithReference(input.Reference).
		WithNeighborhood(input.Neighborhood).
		WithCity(input.City).
		WithState(input.State).
		WithCountry(input.Country).
		WithUpdatedAt(time.Now()).
		Build()
}

func (m *CompanyMapper) PhoneInputToDomain(input command.CreateCompanyPhoneCommand) (*phoneDomain.Phone, error) {
	return phoneDomain.NewPhoneBuilder().
		WithNumber(input.Number).
		WithKind(input.Kind).
		WithDepartment(input.Department).
		Build()
}

func (m *CompanyMapper) UpdatePhoneInputToDomain(input command.UpdateCompanyPhoneCommand) (*phoneDomain.Phone, error) {
	return phoneDomain.NewPhoneBuilder().
		WithID(input.IDPhone).
		WithNumber(input.Number).
		WithKind(input.Kind).
		WithDepartment(input.Department).
		WithUpdatedAt(time.Now()).
		Build()
}

func (m *CompanyMapper) EmailInputToDomain(input command.CreateCompanyEmailCommand) (*emailDomain.Email, error) {
	return emailDomain.NewEmailBuilder().
		WithAddress(input.Address).
		Build()
}

func (m *CompanyMapper) UpdateEmailInputToDomain(input command.UpdateCompanyEmailCommand) (*emailDomain.Email, error) {
	return emailDomain.NewEmailBuilder().
		WithID(input.IDEmail).
		WithAddress(input.Address).
		WithUpdatedAt(time.Now()).
		Build()
}

func (m *CompanyMapper) SocialMediaInputToDomain(input []command.CreateCompanySocialMediaCommand) ([]*socialMediaDomain.SocialMedia, error) {
	socialMedia := make([]*socialMediaDomain.SocialMedia, len(input))
	for i, sm := range input {
		dsm, err := socialMediaDomain.NewSocialMediaBuilder().
			WithPlatform(sm.Platform).
			WithURL(sm.URL).
			Build()
		if err != nil {
			return nil, err
		}
		socialMedia[i] = dsm
	}
	return socialMedia, nil
}

func (m *CompanyMapper) UpdateSocialMediaInputToDomain(input []command.UpdateCompanySocialMediaCommand) ([]*socialMediaDomain.SocialMedia, error) {
	socialMedia := make([]*socialMediaDomain.SocialMedia, len(input))
	for i, sm := range input {
		dsm, err := socialMediaDomain.NewSocialMediaBuilder().
			WithID(sm.IDSocialMedia).
			WithPlatform(sm.Platform).
			WithURL(sm.URL).
			WithUpdatedAt(time.Now()).
			Build()
		if err != nil {
			return nil, err
		}
		socialMedia[i] = dsm
	}
	return socialMedia, nil
}
