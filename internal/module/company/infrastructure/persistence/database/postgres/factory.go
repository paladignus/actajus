// Package postgres
package postgres

import (
	"fmt"

	addr "github.com/paladignus/actajus/internal/module/address/domain"
	addrDB "github.com/paladignus/actajus/internal/module/address/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/company/application/repository"
	"github.com/paladignus/actajus/internal/module/company/domain"
	email "github.com/paladignus/actajus/internal/module/email/domain"
	emailDB "github.com/paladignus/actajus/internal/module/email/infrastructure/persistence/database/postgres"
	phone "github.com/paladignus/actajus/internal/module/phone/domain"
	phoneDB "github.com/paladignus/actajus/internal/module/phone/infrastructure/persistence/database/postgres"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/domain"
	socialMediaDB "github.com/paladignus/actajus/internal/module/social_media/infrastructure/persistence/database/postgres"

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

func (f *Factory) Address() addr.AddressRepository {
	return addrDB.NewAddress(f.exec)
}

func (f *Factory) CompanyAddress() domain.CompanyAddressRepository {
	return NewCompanyAddress(f.exec)
}

func (f *Factory) Phone() phone.PhoneRepository {
	return phoneDB.NewPhone(f.exec)
}

func (f *Factory) CompanyPhone() domain.CompanyPhoneRepository {
	return NewCompanyPhone(f.exec)
}

func (f *Factory) Email() email.EmailRepository {
	return emailDB.NewEmail(f.exec)
}

func (f *Factory) CompanyEmail() domain.CompanyEmailRepository {
	return NewCompanyEmail(f.exec)
}

func (f *Factory) SocialMedia() socialMedia.SocialMediaRepository {
	return socialMediaDB.NewSocialMedia(f.exec)
}
