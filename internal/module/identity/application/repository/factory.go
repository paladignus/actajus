// Package repository
package repository

import (
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type Factory interface {
	WithTx(tx uow.Tx) Factory
	User() UserRepository
	PasswordReset() PasswordResetRepository
	EmailVerification() EmailVerificationRepository
	// Company() CompanyRepository
	// Address() addr.AddressRepository
	// CompanyAddress() domain.CompanyAddressRepository
	// Phone() phone.PhoneRepository
	// CompanyPhone() domain.CompanyPhoneRepository
	// Email() email.EmailRepository
	// CompanyEmail() domain.CompanyEmailRepository
	// SocialMedia() socialMedia.SocialMediaRepository
}
