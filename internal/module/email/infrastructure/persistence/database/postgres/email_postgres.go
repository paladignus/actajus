// Package database
package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/module/email/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Email struct {
	db postgres.Executor
}

func NewEmail(db postgres.Executor) Email {
	return Email{
		db,
	}
}

func (e Email) Create(ctx context.Context, email *domain.Email) error {
	query := `INSERT INTO emails (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING idemails`
	var id uint
	if err := e.db.QueryRow(ctx, query,
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
	query := `UPDATE emails SET address = $1, updated_at = $2 WHERE idemails = $3`
	_, err := e.db.Exec(ctx, query,
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
	_, err := e.db.Exec(ctx, query,
		email.UpdatedAt(),
		email.DeletedAt(),
		email.ID())
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (e Email) FindByIDCompany(ctx context.Context, idCompany uint) (*domain.Email, error) {
	query := `
			SELECT
				e.idemails, e.address, e.created_at, e.updated_at
			FROM emails e
			LEFT JOIN company_email ce ON e.idemails = ce.id_emails AND e.deleted_at IS NULL
			WHERE id_companies = $1`
	var (
		id        uint
		address   string
		createdAt time.Time
		updatedAt time.Time
	)
	err := e.db.QueryRow(ctx, query, idCompany).
		Scan(&id, &address, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return nil, sharedDomain.NewFieldError("id_company", "idCompany")
		}
		return nil, err
	}
	return domain.NewEmailBuilder().
		WithID(id).
		WithAddress(address).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
