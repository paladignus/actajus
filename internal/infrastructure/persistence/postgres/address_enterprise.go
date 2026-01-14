// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type AddressEnterprise struct {
	tx PgxPool
}

func NewAddressEnterprise(tx PgxPool) AddressEnterprise {
	return AddressEnterprise{tx}
}

func (a AddressEnterprise) Create(ctx context.Context, idAddress, idEnterprise uint) error {
	sql := `INSERT INTO address_enterprise (id_addresses, id_companies) VALUES ($1, $2);`
	if _, err := a.tx.Exec(ctx, sql, idAddress, idEnterprise); err != nil {
		return fmt.Errorf("database error while saving address enterprise: %w", err)
	}
	return nil
}
