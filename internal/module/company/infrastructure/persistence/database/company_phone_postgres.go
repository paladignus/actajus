// Package database
package database

import (
	"context"
	"time"

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
	now := time.Now()
	query := `INSERT INTO company_phone (id_companies, id_phones, created_at, started_at)
		VALUES ($1, $2, $3, $4)`
	_, err := c.pool.Exec(ctx, query, idCompany, idPhone, now, now)
	return err
}

func (c CompanyPhone) DeleteByIDCompany(ctx context.Context, idCompany uint) error {
	query := `UPDATE company_phone SET ended_at = $1 WHERE id_companies = $2 AND ended_at IS NULL`
	_, err := c.pool.Exec(ctx, query, time.Now(), idCompany)
	return err
}
