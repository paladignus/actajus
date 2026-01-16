// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type CompanyEmail struct {
	pool PgxPool
}

func NewCompanyEmail(pool PgxPool) CompanyEmail {
	return CompanyEmail{pool}
}

func (c CompanyEmail) Create(ctx context.Context, idCompany, idEmail uint) error {
	sql := `INSERT INTO company_email (id_companies, id_emails) VALUES ($1, $2);`
	if _, err := c.pool.Exec(ctx, sql, idCompany, idEmail); err != nil {
		return fmt.Errorf("database error while saving email company: %w", err)
	}
	return nil
}
