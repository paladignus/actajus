// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/domain/entity"
)

type Enterprise struct {
	tx pgx.Tx
}

func NewEnterprise(tx pgx.Tx) Enterprise {
	return Enterprise{tx}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) (id string, err error) {
	sql := `INSERT INTO companies (registered_by, name, trade_name, cnpj) VALUES ($1, $2, $3, $4) returning idcompanies;`
	if err = e.tx.QueryRow(ctx, sql, enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).Scan(&id); err != nil {
		return "", fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return id, nil
}
