// Package database
package database

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
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
	now := time.Now()
	query := `INSERT INTO company_email (id_companies, id_emails, created_at, started_at)
		VALUES ($1, $2, $3, $4)`
	_, err := c.pool.Exec(ctx, query, idCompany, idEmail, now, now)
	return err
}

func (c CompanyEmail) DeleteByIDCompany(ctx context.Context, idCompany uint) error {
	query := `DELETE FROM company_email WHERE id_companies = $1`
	_, err := c.pool.Exec(ctx, query, idCompany)
	return err
}
