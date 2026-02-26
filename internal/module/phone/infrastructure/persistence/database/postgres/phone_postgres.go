// Package database
package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paladignus/actajus/internal/infrastructure/persistence/postgres"
	"github.com/paladignus/actajus/internal/module/phone/domain"
	sharedDomain "github.com/paladignus/actajus/internal/shared/domain"
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

func (p Phone) FindByIDCompany(ctx context.Context, idCompany uint) (*domain.Phone, error) {
	query := `
			SELECT
				p.idphones, p.number, p.kind, p.department, p.created_at, p.updated_at
			FROM phones p
			LEFT JOIN company_phone cp ON p.idphones = cp.id_phones AND cp.ended_at IS NULL
			WHERE id_companies = $1`
	var (
		id         uint
		number     string
		kind       string
		department string
		createdAt  time.Time
		updatedAt  time.Time
	)
	err := p.pool.QueryRow(ctx, query, idCompany).
		Scan(&id, &number, &kind, &department, &createdAt, &updatedAt)
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
	return domain.NewPhoneBuilder().
		WithID(id).
		WithNumber(number).
		WithKind(kind).
		WithDepartment(department).
		WithCreatedAt(createdAt).
		WithUpdatedAt(updatedAt).
		Build()
}
