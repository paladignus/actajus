// Package database
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	addressDomain "github.com/paladignus/actajus/internal/module/address/domain"
	addressDatabase "github.com/paladignus/actajus/internal/module/address/infrastructure/persistence/database"
	"github.com/paladignus/actajus/internal/module/company/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type CompanyUnitOfWork struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

func NewCompanyUnitOfWork(db *pgxpool.Pool) CompanyUnitOfWork {
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

func (u *CompanyUnitOfWork) Address() addressDomain.AddressRepository {
	return addressDatabase.NewAddress(u.GetPgxPool())
}
