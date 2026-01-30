// Package mapper
package mapper

import (
	"time"

	addrDTO "github.com/paladignus/actajus/internal/module/address/application/dto"
	addrMapper "github.com/paladignus/actajus/internal/module/address/application/mapper"
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/application/dto"
	"github.com/paladignus/actajus/internal/module/company/domain"
	emailDTO "github.com/paladignus/actajus/internal/module/email/application/dto"
	emailMapper "github.com/paladignus/actajus/internal/module/email/application/mapper"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	phoneDTO "github.com/paladignus/actajus/internal/module/phone/application/dto"
	phoneMapper "github.com/paladignus/actajus/internal/module/phone/application/mapper"
	phoneDomain "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMediaDTO "github.com/paladignus/actajus/internal/module/social_media/application/dto"
	socialMediaMapper "github.com/paladignus/actajus/internal/module/social_media/application/mapper"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
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

func (m *CompanyMapper) CompanyInputToDomain(input dto.CreateCompanyRequest) (*domain.Company, error) {
	v := validation.New(validation.PT)
	if err := v.ValidateStruct(input); err != nil {
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithRegisteredBy(input.RegisteredBy).
		Build()
}

func (m *CompanyMapper) UpdateInputDomain(input dto.UpdateCompanyRequest) (*domain.Company, error) {
	v := validation.New(validation.PT)
	if err := v.ValidateStruct(input); err != nil {
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithID(input.IDCompany).
		WithName(input.Name).
		WithTradeName(input.TradeName).
		WithCNPJ(input.CNPJ).
		WithUpdatedAt(time.Now()).
		Build()
}

func (m *CompanyMapper) AddressInputToDomain(input addrDTO.CreateAddressRequest) (*addrDomain.Address, error) {
	return m.addrMapper.InputToDomain(input)
}

func (m *CompanyMapper) UpdateAddressInputToDomain(input addrDTO.UpdateAddressRequest) (*addrDomain.Address, error) {
	return m.addrMapper.UpdateInputToDomain(input)
}

func (m *CompanyMapper) PhoneInputToDomain(input phoneDTO.CreatePhoneRequest) (*phoneDomain.Phone, error) {
	return m.phoneMapper.InputToDomain(input)
}

func (m *CompanyMapper) EmailInputToDomain(input emailDTO.CreateEmailRequest) (*emailDomain.Email, error) {
	return m.emailMapper.InputToDomain(input)
}

func (m *CompanyMapper) SocialMediaInputToDomain(input []socialMediaDTO.CreateSocialMediaRequest) ([]*socialMediaDomain.SocialMedia, error) {
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
