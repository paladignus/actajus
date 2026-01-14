// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Enterprise struct {
	tx PgxPool
}

func NewEnterprise(tx PgxPool) Enterprise {
	return Enterprise{tx}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) (id uint, err error) {
	sql := `INSERT INTO companies (registered_by, name, trade_name, cnpj) VALUES ($1, $2, $3, $4) returning idcompanies;`
	if err = e.tx.QueryRow(ctx, sql, enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).Scan(&id); err != nil {
		return 0, fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return id, nil
}
