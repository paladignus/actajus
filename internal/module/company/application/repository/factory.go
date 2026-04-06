// Package repository
package repository

import (
	"github.com/paladignus/actajus/internal/module/company/domain"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type Factory interface {
	WithTx(tx uow.Tx) Factory
	Company() CompanyRepository
	Address() AddressRepository
	CompanyAddress() domain.CompanyAddressRepository
	Phone() PhoneRepository
	CompanyPhone() domain.CompanyPhoneRepository
	Email() EmailRepository
	CompanyEmail() domain.CompanyEmailRepository
	SocialMedia() SocialMediaRepository
}
