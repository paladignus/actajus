// Package postgres
package postgres

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CompanyAddress struct {
	db postgres.Executor
}

func NewCompanyAddress(db postgres.Executor) *CompanyAddress {
	return &CompanyAddress{
		db,
	}
}

func (c CompanyAddress) Create(ctx context.Context, idCompany, idAddress int64) error {
	query := `INSERT INTO company_address (id_companies, id_addresses, started_at)
		VALUES ($1, $2, $3)`
	_, err := c.db.Exec(ctx, query, idCompany, idAddress, time.Now())
	return err
}

func (c CompanyAddress) DeleteByIDCompany(ctx context.Context, idCompany int64) error {
	query := `UPDATE company_address SET ended_at = $1 WHERE id_companies = $2 AND ended_at IS NULL`
	_, err := c.db.Exec(ctx, query, time.Now(), idCompany)
	return err
}
