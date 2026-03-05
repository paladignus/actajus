// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	addr "github.com/paladignus/actajus/internal/module/address/domain"
	addrDB "github.com/paladignus/actajus/internal/module/address/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/module/company/domain"
	email "github.com/paladignus/actajus/internal/module/email/domain"
	emailDB "github.com/paladignus/actajus/internal/module/email/infrastructure/persistence/database/postgres"
	phone "github.com/paladignus/actajus/internal/module/phone/domain"
	phoneDB "github.com/paladignus/actajus/internal/module/phone/infrastructure/persistence/database/postgres"
	socialMedia "github.com/paladignus/actajus/internal/module/social_media/domain"
	socialMediaDB "github.com/paladignus/actajus/internal/module/social_media/infrastructure/persistence/database/postgres"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CompanyUnitOfWork struct {
	db postgres.PgxPool
	tx pgx.Tx
}

func NewCompanyUnitOfWork(db postgres.PgxPool) CompanyUnitOfWork {
	return CompanyUnitOfWork{db: db}
}

func (u *CompanyUnitOfWork) Begin(ctx context.Context) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	u.tx = tx
	return nil
}

func (u *CompanyUnitOfWork) Commit(ctx context.Context) error {
	if err := u.tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u *CompanyUnitOfWork) Rollback(ctx context.Context) error {
	if err := u.tx.Rollback(ctx); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (u *CompanyUnitOfWork) GetPgxPool() postgres.PgxPool {
	if u.tx != nil {
		return postgres.NewTxAdapter(u.tx)
	}
	return u.db
}

func (u *CompanyUnitOfWork) Company() domain.CompanyRepository {
	return NewCompany(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) Address() addr.AddressRepository {
	return addrDB.NewAddress(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) CompanyAddress() domain.CompanyAddressRepository {
	return NewCompanyAddress(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) Phone() phone.PhoneRepository {
	return phoneDB.NewPhone(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) CompanyPhone() domain.CompanyPhoneRepository {
	return NewCompanyPhone(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) Email() email.EmailRepository {
	return emailDB.NewEmail(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) CompanyEmail() domain.CompanyEmailRepository {
	return NewCompanyEmail(u.GetPgxPool())
}

func (u *CompanyUnitOfWork) SocialMedia() socialMedia.SocialMediaRepository {
	return socialMediaDB.NewSocialMedia(u.GetPgxPool())
}
