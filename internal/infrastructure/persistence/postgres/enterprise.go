// Package postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/paladignus/actajus/internal/domain/entity"
	"github.com/paladignus/actajus/internal/infrastructure/database"
)

type Enterprise struct {
	db database.PgxPool
}

func NewEnterprise(db database.PgxPool) *Enterprise {
	return &Enterprise{db: db}
}

func (e Enterprise) Create(ctx context.Context, enterprise entity.Enterprise) (id string, err error) {
	sql := `INSERT INTO companies (registered_by, name, trade_name, cnpj) VALUES ($1, $2, $3, $4) returning idcompanies;`
	err = e.db.QueryRow(ctx, sql, enterprise.RegisteredBy, enterprise.Name, enterprise.TradeName, enterprise.CNPJ).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("database error while saving enterprise: %w", err)
	}
	return id, nil
}
