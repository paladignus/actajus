// Package database
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/postgres"
)

type CompanyAddress struct {
	pool postgres.PgxPool
}

func NewCompanyAddress(pool postgres.PgxPool) *CompanyAddress {
	return &CompanyAddress{
		pool,
	}
}

func (c CompanyAddress) Create(ctx context.Context, idCompany, idAddress uint) error {
	now := time.Now()
	fmt.Println(idCompany, idAddress, now)
	query := `INSERT INTO company_address (id_companies, id_addresses, created_at, started_at)
		VALUES ($1, $2, $3, $4)`
	_, err := c.pool.Exec(ctx, query, idCompany, idAddress, now, now)
	return err
}
