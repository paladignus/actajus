// Package postgres
package postgres

import (
	"context"
	"fmt"
)

type EmailEnterprise struct {
	db PgxPool
}

func NewEmailEnterprise(db PgxPool) EmailEnterprise {
	return EmailEnterprise{db}
}

func (e EmailEnterprise) Create(ctx context.Context, idEmail, idEnterprise uint) error {
	sql := `INSERT INTO email_enterprise (id_emails, id_companies) VALUES ($1, $2);`
	if _, err := e.db.Exec(ctx, sql, idEmail, idEnterprise); err != nil {
		return fmt.Errorf("database error while saving email enterprise: %w", err)
	}
	return nil
}
