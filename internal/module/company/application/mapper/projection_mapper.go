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
)

type CompanyProjectionMapper struct {
	addrMapper  addrMapper.AddressPrejectionMapper
	phoneMapper phoneMapper.PhoneProjectionMapper
	emailMapper emailMapper.EmailProjectionMapper
}

func NewCompanyProjectionMapper(
	addrMapper addrMapper.AddressPrejectionMapper,
	phoneMapper phoneMapper.PhoneProjectionMapper,
	emailMapper emailMapper.EmailProjectionMapper,
) CompanyProjectionMapper {
	return CompanyProjectionMapper{
		addrMapper,
		phoneMapper,
		emailMapper,
	}
}

func (m *CompanyProjectionMapper) ProjectCompanyToReadModel(
	company *domain.Company,
	address *addrDomain.Address,
	phone *phoneDomain.Phone,
	email *emailDomain.Email,
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
	return readModel
}
