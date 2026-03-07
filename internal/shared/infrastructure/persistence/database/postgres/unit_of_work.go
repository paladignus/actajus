// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(tx uow.Tx) error) error {
	tx, err := u.pool.BeginTx(ctx, pgx.TxOptions{})
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

// // SimpleUnitOfWork implementa unitofwork.UnitOfWork usando PgxPool
// // Este é um UnitOfWork simples para operações que não precisam de múltiplos repositórios
// type SimpleUnitOfWork struct {
// 	pool PgxPool
// 	tx   pgx.Tx
// }
//
// // NewSimpleUnitOfWork cria uma nova instância de SimpleUnitOfWork
// func NewSimpleUnitOfWork(pool PgxPool) *SimpleUnitOfWork {
// 	return &SimpleUnitOfWork{pool: pool}
// }
//
// // Begin inicia uma nova transação
// func (u *SimpleUnitOfWork) Begin(ctx context.Context) error {
// 	tx, err := u.pool.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	u.tx = tx
// 	return nil
// }
//
// // Commit confirma a transação
// func (u *SimpleUnitOfWork) Commit(ctx context.Context) error {
// 	if u.tx == nil {
// 		return fmt.Errorf("no active transaction to commit")
// 	}
// 	if err := u.tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	u.tx = nil
// 	return nil
// }
//
// // Rollback desfaz a transação
// func (u *SimpleUnitOfWork) Rollback(ctx context.Context) error {
// 	if u.tx == nil {
// 		return nil // No transaction to rollback
// 	}
// 	if err := u.tx.Rollback(ctx); err != nil {
// 		return fmt.Errorf("failed to rollback transaction: %w", err)
// 	}
// 	u.tx = nil
// 	return nil
// }
//
// // Tx retorna a transação atual (para uso interno dos repositórios)
// func (u *SimpleUnitOfWork) Tx() pgx.Tx {
// 	return u.tx
// }
//
// // PgxPool é o pool de conexões (já definido em connection.go)
// // Esta struct apenas usa a interface existente
//
