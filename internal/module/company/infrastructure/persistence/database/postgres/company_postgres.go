// Package postgres
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/module/company/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type Company struct {
	db postgres.Executor
}

func NewCompany(db postgres.Executor) Company {
	return Company{
		db,
	}
}

func (c Company) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (registered_by, name, trade_name, cnpj, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING idcompanies`
	var id int64
	if err := c.db.QueryRow(ctx, query,
		company.RegisteredBy(),
		company.Name().Value(),
		company.TradeName().Value(),
		company.CNPJ().Value(),
		company.CreatedAt(),
		company.UpdatedAt(),
	).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sharedDomain.NewFieldError("cnpj", "cnpj already exists")
		}
		return err
	}
	return company.SetID(id)
}

func (c Company) Update(ctx context.Context, company *domain.Company) error {
	const query = `
		UPDATE companies
		SET name = $1, trade_name = $2, cnpj = $3, updated_at = $4 
		WHERE idcompanies = $5 AND deleted_at IS NULL`
	_, err := c.db.Exec(ctx, query,
		company.Name().Value(),
		company.TradeName().Value(),
		company.CNPJ().Value(),
		company.UpdatedAt(),
		company.ID(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sharedDomain.NewFieldError("cnpj", "cnpj already exists")
	}
	return err
}

func (c Company) Delete(ctx context.Context, company *domain.Company) error {
	const query = `UPDATE companies SET updated_at = $1, deleted_at = $2 WHERE idcompanies = $3 AND deleted_at IS NULL`
	_, err := c.db.Exec(ctx, query,
		company.UpdatedAt(),
		company.DeletedAt(),
		company.ID(),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return sharedDomain.NewFieldError("deleted_at", "company already deleted")
	}
	return err
}

func (c Company) FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	const query = `
		SELECT
			idcompanies, registered_by, name, trade_name, created_at, updated_at
		FROM companies
		WHERE cnpj = $1 AND deleted_at IS NULL`
	var (
		id           int64
		registeredBy int64
		name         string
		tradeName    string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := c.db.QueryRow(ctx, query, cnpj).
		Scan(
			&id, &registeredBy, &name, &tradeName, &createdAt, &updatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedDomain.NewFieldError("cnpj", "cnpj not found")
		}
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithID(id).
		WithRegisteredBy(registeredBy).
		WithName(name).
		WithTradeName(tradeName).
		WithCNPJ(cnpj).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}

func (c Company) FindByID(ctx context.Context, id int64) (*domain.Company, error) {
	const query = `
			SELECT
					registered_by, name, trade_name, cnpj, created_at, updated_at
			FROM companies
			WHERE idcompanies = $1 AND deleted_at IS NULL`
	var (
		registeredBy int64
		name         string
		tradeName    string
		cnpj         string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := c.db.QueryRow(ctx, query, id).
		Scan(
			&registeredBy, &name, &tradeName, &cnpj, &createdAt, &updatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharedDomain.NewFieldError("id", "company not found")
		}
		return nil, err
	}
	return domain.NewCompanyBuilder().
		WithID(id).
		WithRegisteredBy(registeredBy).
		WithName(name).
		WithTradeName(tradeName).
		WithCNPJ(cnpj).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
