// Package database
package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/phone/domain"
)

type Phone struct {
	pool postgres.PgxPool
}

func NewPhone(pool postgres.PgxPool) Phone {
	return Phone{
		pool,
	}
}

func (p Phone) Create(ctx context.Context, phone *domain.Phone) error {
	query := `INSERT INTO phones (number, kind, department, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING idphones`
	var id uint
	if err := p.pool.QueryRow(ctx, query,
		phone.Number(),
		phone.Kind(),
		phone.Department(),
		phone.CreatedAt(),
		phone.UpdatedAt()).Scan(&id); err != nil {
		return err
	}
	return phone.SetID(id)
}

func (p Phone) Update(ctx context.Context, phone domain.Phone) error {
	query := `UPDATE phones SET number = $1, kind = $2, department = $3, updated_at = $4 WHERE idphones = $5`
	_, err := p.pool.Exec(ctx, query,
		phone.Number(),
		phone.Kind(),
		phone.Department(),
		phone.UpdatedAt(),
		phone.ID())
	return err
}

func (p Phone) Delete(ctx context.Context, phone domain.Phone) error {
	query := `UPDATE phones SET updated_at = $1, deleted_at = $2 WHERE idphones = $3 AND deleted_at IS NULL`
	_, err := p.pool.Exec(ctx, query,
		phone.UpdatedAt(),
		phone.DeletedAt(),
		phone.ID())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}
