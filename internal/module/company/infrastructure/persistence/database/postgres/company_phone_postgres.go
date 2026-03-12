// Package postgres
package postgres

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type CompanyPhone struct {
	db postgres.Executor
}

func NewCompanyPhone(db postgres.Executor) *CompanyPhone {
	return &CompanyPhone{
		db,
	}
}

func (c CompanyPhone) Create(ctx context.Context, idCompany, idPhone uint) error {
	query := `INSERT INTO company_phone (id_companies, id_phones, started_at)
		VALUES ($1, $2, $3)`
	_, err := c.db.Exec(ctx, query, idCompany, idPhone, time.Now())
	return err
}

func (c CompanyPhone) DeleteByIDCompany(ctx context.Context, idCompany uint) error {
	query := `UPDATE company_phone SET ended_at = $1 WHERE id_companies = $2 AND ended_at IS NULL`
	_, err := c.db.Exec(ctx, query, time.Now(), idCompany)
	return err
}
