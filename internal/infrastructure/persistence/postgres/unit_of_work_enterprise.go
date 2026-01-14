// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/domain/repository"
)

type UnitOfWorkEnterprise struct {
	db *pgxpool.Pool
	tx pgx.Tx
}

func NewUnitOfWorkEnterprise(db *pgxpool.Pool) *UnitOfWorkEnterprise {
	return &UnitOfWorkEnterprise{db: db}
}

func (u *UnitOfWorkEnterprise) Begin(ctx context.Context) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	u.tx = tx
	return nil
}

func (u *UnitOfWorkEnterprise) Commit(ctx context.Context) error {
	if err := u.tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (u *UnitOfWorkEnterprise) Rollback(ctx context.Context) error {
	if err := u.tx.Rollback(ctx); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (u *UnitOfWorkEnterprise) Enterprise() repository.IEnterprise {
	return NewEnterprise(u.tx)
}

func (u *UnitOfWorkEnterprise) Address() repository.IAddress {
	return NewAddress(u.tx)
}

func (u *UnitOfWorkEnterprise) AddressEnterprise() repository.IAddressEnterprise {
	return NewAddressEnterprise(u.tx)
}
