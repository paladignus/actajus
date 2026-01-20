// Package domain
package domain

import (
	address "github.com/paladignus/actajus/internal/module/address/domain"
	uow "github.com/paladignus/actajus/internal/shared/domain/unit_of_work"
)

type CompanyUnitOfWork interface {
	uow.UnitOfWork
	Company() CompanyRepository
	Address() address.AddressRepository
	CompanyAddress() CompanyAddressRepository
}
