// Package repository
package repository

import (
	addr "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/domain"
	email "github.com/paladignus/actajus/internal/module/email/domain"
	phone "github.com/paladignus/actajus/internal/module/phone/domain"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/domain"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type Factory interface {
	WithTx(tx uow.Tx) Factory
	Company() CompanyRepository
	Address() addr.AddressRepository
	CompanyAddress() domain.CompanyAddressRepository
	Phone() phone.PhoneRepository
	CompanyPhone() domain.CompanyPhoneRepository
	Email() email.EmailRepository
	CompanyEmail() domain.CompanyEmailRepository
	SocialMedia() socialMedia.SocialMediaRepository
}
