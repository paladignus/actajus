// Package postgres
package postgres

import (
	"context"

	"github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
)

type SentEmailRepository struct {
	db postgres.Executor
}

func NewSentEmailRepository(db postgres.Executor) *SentEmailRepository {
	return &SentEmailRepository{db}
}

func (r *SentEmailRepository) ExistsByIDMessage(ctx context.Context, mid string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM sent_emails WHERE id_message = $1)`,
		mid).Scan(&exists)
	return exists, err
}

func (r *SentEmailRepository) MarkSent(ctx context.Context, mid string, recipient string, template string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sent_emails (id_message, recipient, template)
		VALUES ($1, $2, $3)
		ON CONFLICT (id_message) DO NOTHING
	`, mid, recipient, template)
	return err
}
