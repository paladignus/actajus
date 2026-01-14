// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type EnterprisePhone struct {
	tx PgxPool
}

func NewEnterprisePhone(tx PgxPool) EnterprisePhone {
	return EnterprisePhone{tx}
}

func (e EnterprisePhone) Create(ctx context.Context, idEnterprise, idPhone uint) error {
	sql := `INSERT INTO enterprise_phone (id_companies, id_phones) VALUES ($1, $2);`
	if _, err := e.tx.Exec(ctx, sql, idEnterprise, idPhone); err != nil {
		return fmt.Errorf("database error while saving enterprise phone: %w", err)
	}
	return nil
}
