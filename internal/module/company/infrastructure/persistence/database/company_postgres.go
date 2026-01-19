// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/module/company/domain"
	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type Company struct {
	pool postgres.PgxPool
}

func NewCompany(pool postgres.PgxPool) *Company {
	return &Company{
		pool,
	}
}

func (c *Company) FindByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	return nil, nil
}

func (c Company) Create(ctx context.Context, company *domain.Company) error {
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
