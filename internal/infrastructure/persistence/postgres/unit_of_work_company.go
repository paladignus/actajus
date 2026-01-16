// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type UnitOfWorkCompany struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

func NewUnitOfWorkCompany(db *pgxpool.Pool) *UnitOfWorkCompany {
	return &UnitOfWorkCompany{db: db}
}

func (u *UnitOfWorkCompany) Begin(ctx context.Context) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	u.tx = tx
	return nil
}

func (u *UnitOfWorkCompany) Commit(ctx context.Context) error {
	if err := u.tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u *UnitOfWorkCompany) Rollback(ctx context.Context) error {
	if err := u.tx.Rollback(ctx); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (u *UnitOfWorkCompany) GetPgxPool() PgxPool {
	if u.tx != nil {
		return NewTxAdapter(u.tx)
	}
	return u.db
}

func (u *UnitOfWorkCompany) Company() repository.ICompany {
	return NewCompany(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) Address() repository.IAddress {
	return NewAddress(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) CompanyAddress() repository.ICompanyAddress {
	return NewCompanyAddress(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) Phone() repository.IPhone {
	return NewPhone(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) CompanyPhone() repository.ICompanyPhone {
	return NewCompanyPhone(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) Email() repository.IEmail {
	return NewEmail(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) CompanyEmail() repository.ICompanyEmail {
	return NewCompanyEmail(u.GetPgxPool())
}

func (u *UnitOfWorkCompany) SocialMedia() repository.ISocialMedia {
	return NewSocialMedia(u.GetPgxPool())
}
