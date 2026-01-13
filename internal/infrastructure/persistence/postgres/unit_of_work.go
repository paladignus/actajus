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
	pool       *pgxpool.Pool
	tx         pgx.Tx
	completed  bool
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
	uow.completed = false // Reset completion flag for new transaction
	return nil
}

func (uow *UnitOfWork) Commit(ctx context.Context) error {
	if uow.tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}
	if uow.completed {
		return fmt.Errorf("transaction already completed")
	}
	err := uow.tx.Commit(ctx)
	if err == nil {
		uow.completed = true
	}
	return err
}

func (uow *UnitOfWork) Rollback(ctx context.Context) error {
	if uow.tx == nil {
		return fmt.Errorf("no active transaction to rollback")
	}
	if uow.completed {
		return fmt.Errorf("transaction already completed")
	}
	err := uow.tx.Rollback(ctx)
	if err == nil {
		uow.completed = true
	}
	return err
}

func (uow *UnitOfWork) GetPgxPool() PgxPool {
	if uow.tx != nil {
		return NewTxAdapter(uow.tx)
	}
	return uow.pool
}
