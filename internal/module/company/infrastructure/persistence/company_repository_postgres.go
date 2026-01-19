// Package persistence
package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewCompanyRepositoryPostgres(pool *pgxpool.Pool) *CompanyRepositoryPostgres {
	return &CompanyRepositoryPostgres{
		pool,
	}
}

func (c *CompanyRepositoryPostgres) FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	return nil, nil
}

func (c CompanyRepositoryPostgres) Create(ctx context.Context, company *domain.Company) error {
	query := `INSERT INTO companies (registered_by, name, trade_name, cnpj, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING idcompanies`
	var id int
	if err := c.pool.QueryRow(ctx, query,
		company.RegisteredBy(),
		company.Name().Value(),
		company.TradeName().Value(),
		company.CNPJ().Value(),
		company.CreatedAt(),
		company.UpdatedAt(),
	).Scan(&id); err != nil {
		return err
	}
	return company.SetID(id)
}
