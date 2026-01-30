// Package database
package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/email/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
)

type Email struct {
	pool postgres.PgxPool
}

func NewEmail(pool postgres.PgxPool) Email {
	return Email{
		pool,
	}
}

func (e Email) Create(ctx context.Context, email domain.Email) error {
	query := `INSERT INTO emails (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING idemails`
	var id uint
	if err := e.pool.QueryRow(ctx, query,
		email.Address(),
		email.CreatedAt(),
		email.UpdatedAt()).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sharedDomain.NewFieldError("email", "email already exists")
		}
		return err
	}
	return email.SetID(id)
}

func (e Email) Update(ctx context.Context, email domain.Email) error {
	query := `UPDATE emails SET address = $1, updated_at = $2 WHERE idemails = $3 AND deleted_at IS NULL`
	_, err := e.pool.Exec(ctx, query,
		email.Address(),
		email.UpdatedAt(),
		email.ID())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("email", "email already exists")
	}
	return err
}

func (e Email) Delete(ctx context.Context, email domain.Email) error {
	query := `UPDATE emails SET updated_at = $1, deleted_at = $2 WHERE idemails = $3 AND deleted_at IS NULL`
	_, err := e.pool.Exec(ctx, query,
		email.UpdatedAt(),
		email.DeletedAt(),
		email.ID())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}
