// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type CompanyAddress struct {
	pool PgxPool
}

func NewCompanyAddress(pool PgxPool) CompanyAddress {
	return CompanyAddress{pool}
}

func (c CompanyAddress) Create(ctx context.Context, idCompany, idAddress uint) error {
	sql := `INSERT INTO company_address (id_companies, id_addresses) VALUES ($1, $2);`
	if _, err := c.pool.Exec(ctx, sql, idCompany, idAddress); err != nil {
		return fmt.Errorf("database error while saving address company: %w", err)
	}
	return nil
}
