// Package database
package database

import (
	"context"

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
	query := `INSERT INTO company_email (id_companies, id_emails)
		VALUES ($1, $2)`
	_, err := c.pool.Exec(ctx, query, idCompany, idEmail)
	return err
}
