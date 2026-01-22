// Package database
package database

import (
	"context"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type CompanyPhone struct {
	pool postgres.PgxPool
}

func NewCompanyPhone(pool postgres.PgxPool) *CompanyPhone {
	return &CompanyPhone{
		pool,
	}
}

func (c CompanyPhone) Create(ctx context.Context, idCompany, idPhone uint) error {
	query := `INSERT INTO company_phone (id_companies, id_phones)
		VALUES ($1, $2)`
	_, err := c.pool.Exec(ctx, query, idCompany, idPhone)
	return err
}
