// Package mapper
package mapper

import (
	"time"

	addrCommand "github.com/paladignus/actajus/internal/module/address/application/command"
	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/command"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailCommand "github.com/paladignus/actajus/internal/module/email/application/command"
	emailMapper "github.com/paladignus/actajus/internal/module/email/application/mapper"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneCommand "github.com/paladignus/actajus/internal/module/phone/application/command"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaCommand "github.com/paladignus/actajus/internal/module/social_media/application/command"
	socialMediaMapper "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
)

type CompanyMapper struct {
	addrMapper        *addrMapper.AddressMapper
	phoneMapper       *phoneMapper.PhoneMapper
	emailMapper       *emailMapper.EmailMapper
	socialMediaMapper *socialMediaMapper.SocialMediaMapper
}

func NewCompanyMapper(
	addrMapper *addrMapper.AddressMapper,
	phoneMapper *phoneMapper.PhoneMapper,
	emailMapper *emailMapper.EmailMapper,
	socialMediaMapper *socialMediaMapper.SocialMediaMapper,
) *CompanyMapper {
	return &CompanyMapper{
		addrMapper,
		phoneMapper,
		emailMapper,
		socialMediaMapper,
	}
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

func (m *CompanyMapper) AddressInputToDomain(input addrCommand.CreateAddressCommand) (*addrDomain.Address, error) {
	return m.addrMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateAddressInputToDomain(input addrCommand.UpdateAddressCommand) (*addrDomain.Address, error) {
	return m.addrMapper.UpdateInputToDomain(input)
}

func (m *CompanyMapper) PhoneInputToDomain(input phoneCommand.CreatePhoneCommand) (*phoneDomain.Phone, error) {
	return m.phoneMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdatePhoneInputToDomain(input phoneCommand.UpdatePhoneCommand) (*phoneDomain.Phone, error) {
	return m.phoneMapper.UpdateInputToDomain(input)
}

func (m *CompanyMapper) EmailInputToDomain(input emailCommand.CreateEmailCommand) (*emailDomain.Email, error) {
	return m.emailMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateEmailInputToDomain(input emailCommand.UpdateEmailCommand) (*emailDomain.Email, error) {
	return m.emailMapper.UpdateInputToDomain(input)
}

func (m *CompanyMapper) SocialMediaInputToDomain(input []socialMediaCommand.CreateSocialMediaCommand) ([]*socialMediaDomain.SocialMedia, error) {
	socialMedia := make([]*socialMediaDomain.SocialMedia, len(input))
	for i, sm := range input {
		dsm, err := m.socialMediaMapper.InputToDomain(sm)
		if err != nil {
			return nil, err
		}
		socialMedia[i] = dsm
	}
	return socialMedia, nil
}

func (m *CompanyMapper) UpdateSocialMediaInputToDomain(input []socialMediaCommand.UpdateSocialMediaCommand) ([]*socialMediaDomain.SocialMedia, error) {
	socialMedia := make([]*socialMediaDomain.SocialMedia, len(input))
	for i, sm := range input {
		dsm, err := m.socialMediaMapper.UpdateInputToDomain(sm)
		if err != nil {
			return nil, err
		}
		socialMedia[i] = dsm
	}
	return socialMedia, nil
}
