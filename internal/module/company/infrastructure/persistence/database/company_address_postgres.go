// Package database
package database

import (
	"context"
	"time"

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
	query := `INSERT INTO company_address (id_companies, id_addresses, started_at)
		VALUES ($1, $2, $3)`
	_, err := c.pool.Exec(ctx, query, idCompany, idAddress, time.Now())
	return err
}

func (c CompanyAddress) DeleteByIDCompany(ctx context.Context, idCompany uint) error {
	query := `UPDATE company_address SET ended_at = $1 WHERE id_companies = $2 AND ended_at IS NULL`
	_, err := c.pool.Exec(ctx, query, time.Now(), idCompany)
	return err
}
