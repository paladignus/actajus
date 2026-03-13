package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/paladignus/actajus/internal/shared/application/messaging"
)

type OutboxRepository struct {
	exec Executor
}

func (r *OutboxRepository) Add(ctx context.Context, msg messaging.OutboxMessage) error {
	// payload já vem como []byte JSON, mas vamos garantir que é JSONB no insert
	var raw json.RawMessage = msg.Payload
	_, err := r.exec.Exec(ctx, `
		INSERT INTO outbox_messages (id, subject, payload, occurred_at)
		VALUES ($1, $2, $3, $4)
	`, msg.ID, msg.Subject, raw, msg.OccurredAt)
	if err != nil {
		return fmt.Errorf("insert outbox message: %w", err)
	}
	return nil
}

// func NewOutboxRepository() *OutboxRepository {
// 	return &OutboxRepository{}
// }
//
// func (r *OutboxRepository) Add(ctx context.Context, tx uow.Tx, msg messaging.OutboxMessage) error {
// 	pgxTx, ok := tx.(pgx.Tx)
// 	if !ok {
// 		return fmt.Errorf("invalid tx type")
// 	}
// 	var payload json.RawMessage = msg.Payload
// 	_, err := pgxTx.Exec(ctx, `
// 		INSERT INTO outbox_messages (
// 			id,
// 			subject,
// 			payload,
// 			occurred_at
// 		) VALUES ($1, $2, $3, $4)
// 	`,
// 		msg.ID,
// 		msg.Subject,
// 		payload,
// 		msg.OccurredAt,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert outbox message: %w", err)
// 	}
// 	return nil
// }
