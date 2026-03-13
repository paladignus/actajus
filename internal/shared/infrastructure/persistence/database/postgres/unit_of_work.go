// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type UnitOfWork struct {
	db Transactor
}

func NewUnitOfWork(db Transactor) *UnitOfWork {
	return &UnitOfWork{db}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(tx uow.Tx) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
