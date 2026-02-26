// Package postgres
package postgres

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CompanyEmail struct {
	pool postgres.PgxPool
}

func NewCompanyEmail(pool postgres.PgxPool) *CompanyEmail {
	return &CompanyEmail{
		pool,
	}
}

func (c CompanyEmail) Create(ctx context.Context, idCompany, idEmail uint) error {
	query := `INSERT INTO company_email (id_companies, id_emails, started_at)
		VALUES ($1, $2, $3)`
	_, err := c.pool.Exec(ctx, query, idCompany, idEmail, time.Now())
	return err
}

func (c CompanyEmail) DeleteByIDCompany(ctx context.Context, idCompany uint) error {
	query := `UPDATE company_email SET ended_at = $1 WHERE id_companies = $2 AND ended_at IS NULL`
	_, err := c.pool.Exec(ctx, query, time.Now(), idCompany)
	return err
}
