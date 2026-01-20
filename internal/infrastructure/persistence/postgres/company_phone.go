// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type CompanyPhone struct {
	pool PgxPool
}

func NewCompanyPhone(pool PgxPool) CompanyPhone {
	return CompanyPhone{pool}
}

func (c CompanyPhone) Create(ctx context.Context, idCompany, idPhone uint) error {
	sql := `INSERT INTO company_phone (id_companies, id_phones) VALUES ($1, $2);`
	if _, err := c.pool.Exec(ctx, sql, idCompany, idPhone); err != nil {
		return fmt.Errorf("database error while saving enterprise phone: %w", err)
	}
	return nil
}
