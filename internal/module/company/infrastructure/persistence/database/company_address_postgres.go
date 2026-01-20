// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type CompanyAddress struct {
	pool postgres.PgxPool
}

func NewCompanyAddress(pool postgres.PgxPool) *CompanyAddress {
	return &CompanyAddress{
		pool,
	}
}

func (c CompanyAddress) Create(ctx context.Context, idCompany, idAddress uint) error {
	query := `INSERT INTO company_address (id_companies, id_addresses)
		VALUES ($1, $2)`
	_, err := c.pool.Exec(ctx, query, idCompany, idAddress)
	return err
}
