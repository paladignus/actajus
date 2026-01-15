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

func (e Enterprise) Update(ctx context.Context, enterprise entity.Enterprise) error {
	sql := `UPDATE companies SET name = $1, trade_name = $2, cnpj = $3, updated_at = now() WHERE idcompanies = $4;`
	if _, err := e.tx.Exec(ctx, sql, enterprise.Name, enterprise.TradeName, enterprise.CNPJ, enterprise.IDEnterprise); err != nil {
		return fmt.Errorf("database error while update enterprise: %w", err)
	}
	return nil
}

func (e Enterprise) GetAll(ctx context.Context) (enterprises []entity.Enterprise, err error) {
	sql := `SELECT idcompanies, registered_by, name, trade_name, cnpj FROM companies WHERE deleted_at IS NULL;`
	rows, err := e.tx.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("database error while getting all enterprises: %w", err)
	}
	// enterprises = make([]entity.Enterprise, 0)
	for rows.Next() {
		var enterprise entity.Enterprise
		if err = rows.Scan(&enterprise.IDEnterprise, &enterprise.RegisteredBy, &enterprise.Name, &enterprise.TradeName, &enterprise.CNPJ); err != nil {
			return nil, fmt.Errorf("database error while getting all enterprises: %w", err)
		}
		enterprises = append(enterprises, enterprise)
	}
	return enterprises, nil
}
