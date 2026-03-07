// Package postgres
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SentEmailRepository struct {
	pool *pgxpool.Pool
}

func NewSentEmailRepository(pool *pgxpool.Pool) *SentEmailRepository {
	return &SentEmailRepository{pool: pool}
}

func (r *SentEmailRepository) ExistsByMessageID(ctx context.Context, mid string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM sent_emails
			WHERE id_message = $1
		)
	`, mid).Scan(&exists)

	return exists, err
}

func (r *SentEmailRepository) MarkSent(
	ctx context.Context,
	mid string,
	recipient string,
	template string,
) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sent_emails (id_message, recipient, template)
		VALUES ($1, $2, $3)
		ON CONFLICT (id_message) DO NOTHING
	`, mid, recipient, template)
	return err
}
