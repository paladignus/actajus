// Package postgres
package postgres

import (
	"fmt"

	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/module/company/domain"

	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"

	"github.com/jackc/pgx/v5"
)

type Factory struct {
	exec postgres.Executor
}

func NewFactory(exec postgres.Executor) *Factory {
	return &Factory{exec}
}

func (f *Factory) WithTx(tx uow.Tx) repository.Factory {
	pgxTx, ok := tx.(pgx.Tx)
	if !ok {
		panic(fmt.Sprintf("invalid tx type: %T", tx))
	}
	return &Factory{exec: pgxTx}
}

func (f *Factory) Company() repository.CompanyRepository {
	return NewCompany(f.exec) // seu repo concreto deve aceitar Executor
}

func (f *Factory) Address() repository.AddressRepository {
	return NewAddressRepository(f.exec)
}

func (f *Factory) CompanyAddress() domain.CompanyAddressRepository {
	return NewCompanyAddress(f.exec)
}

func (f *Factory) Phone() repository.PhoneRepository {
	return NewPhoneRepository(f.exec)
}

func (f *Factory) CompanyPhone() domain.CompanyPhoneRepository {
	return NewCompanyPhone(f.exec)
}

func (f *Factory) Email() repository.EmailRepository {
	return NewEmailRepository(f.exec)
}

func (f *Factory) CompanyEmail() domain.CompanyEmailRepository {
	return NewCompanyEmail(f.exec)
}

func (f *Factory) SocialMedia() repository.SocialMediaRepository {
	return NewSocialMediaRepository(f.exec)
}
