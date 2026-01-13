// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
)

type Enterprise struct {
	db PgxPool
}

func NewEnterprise(db PgxPool) *Enterprise {
	return &Enterprise{db: db}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) (id string, err error) {
	sql := `INSERT INTO companies (registered_by, name, trade_name, cnpj) VALUES ($1, $2, $3, $4) returning idcompanies;`
	if err = e.db.QueryRow(ctx, sql, enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).Scan(&id); err != nil {
		return "", fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return id, nil
}
