// Package domain
package domain

import (
	address "github.com/paladignus/actajus/internal/module/address/domain"
	email "github.com/paladignus/actajus/internal/module/email/domain"
	phone "github.com/paladignus/actajus/internal/module/phone/domain"
	uow "github.com/paladignus/actajus/internal/shared/domain/unit_of_work"
)

type CompanyUnitOfWork interface {
	uow.UnitOfWork
	Company() CompanyRepository
	Address() address.AddressRepository
	CompanyAddress() CompanyAddressRepository
	Phone() phone.PhoneRepository
	CompanyPhone() CompanyPhoneRepository
	Email() email.EmailRepository
	CompanyEmail() CompanyEmailRepository
}
