package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type OutboxRepository struct{}

func NewOutboxRepository() *OutboxRepository {
	return &OutboxRepository{}
}

func (r *OutboxRepository) Add(ctx context.Context, tx uow.Tx, msg messaging.OutboxMessage) error {
	pgxTx, ok := tx.(pgx.Tx)
	if !ok {
		return fmt.Errorf("invalid tx type")
	}
	var payload json.RawMessage = msg.Payload
	_, err := pgxTx.Exec(ctx, `
		INSERT INTO outbox_messages (
			id,
			subject,
			payload,
			occurred_at
		) VALUES ($1, $2, $3, $4)
	`,
		msg.ID,
		msg.Subject,
		payload,
		msg.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("insert outbox message: %w", err)
	}
	return nil
}
