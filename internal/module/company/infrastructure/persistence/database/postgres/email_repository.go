// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	emailDomain "github.com/paladignus/actajus/internal/module/email/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type EmailRepository struct {
	exec sharedpostgres.Executor
}

func NewEmailRepository(exec sharedpostgres.Executor) EmailRepository {
	return EmailRepository{exec: exec}
}

func (r EmailRepository) Create(ctx context.Context, email *emailDomain.Email) error {
	query := `INSERT INTO emails (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING idemails`
	var id int64
	if err := r.exec.QueryRow(ctx, query,
		email.Address(),
		email.CreatedAt(),
		email.UpdatedAt(),
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sharedDomain.NewFieldError("email", "email already exists")
		}
		return err
	}
	return email.SetID(id)
}

func (r EmailRepository) Update(ctx context.Context, email emailDomain.Email) error {
	query := `UPDATE emails SET address = $1, updated_at = $2 WHERE idemails = $3`
	_, err := r.exec.Exec(ctx, query,
		email.Address(),
		email.UpdatedAt(),
		email.ID().Value(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("email", "email already exists")
	}
	return err
}

func (r EmailRepository) Delete(ctx context.Context, email emailDomain.Email) error {
	query := `UPDATE emails SET updated_at = $1, deleted_at = $2 WHERE idemails = $3 AND deleted_at IS NULL`
	_, err := r.exec.Exec(ctx, query,
		email.UpdatedAt(),
		email.DeletedAt(),
		email.ID().Value(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil
	}
	return err
}

func (r EmailRepository) FindByIDCompany(ctx context.Context, idCompany int64) (*emailDomain.Email, error) {
	query := `
		SELECT
			e.idemails, e.address, e.created_at, e.updated_at
		FROM emails e
		LEFT JOIN company_email ce ON e.idemails = ce.id_emails AND e.deleted_at IS NULL
		WHERE id_companies = $1`
	var (
		id        int64
		address   string
		createdAt time.Time
		updatedAt time.Time
	)
	err := r.exec.QueryRow(ctx, query, idCompany).
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
	return emailDomain.NewEmailBuilder().
		WithID(id).
		WithAddress(address).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
