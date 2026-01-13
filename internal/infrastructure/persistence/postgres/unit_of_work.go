// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUnitOfWork interface {
	Begin(context.Context) error
	Commit(context.Context) error
	Rollback(context.Context) error
	GetPgxPool() PgxPool
}

type UnitOfWork struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (uow *UnitOfWork) Begin(ctx context.Context) error {
	tx, err := uow.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	uow.tx = tx
	return nil
}

func (uow *UnitOfWork) Commit(ctx context.Context) error {
	if uow.tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}
	return uow.tx.Commit(ctx)
}

func (uow *UnitOfWork) Rollback(ctx context.Context) error {
	if uow.tx == nil {
		return fmt.Errorf("no active transaction to rollback")
	}
	return uow.tx.Rollback(ctx)
}

func (uow *UnitOfWork) GetPgxPool() PgxPool {
	if uow.tx != nil {
		return NewTxAdapter(uow.tx)
	}
	return uow.pool
}
