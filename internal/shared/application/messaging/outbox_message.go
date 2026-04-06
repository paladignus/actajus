// Package messaging
package messaging

import (
	"context"
	"time"

	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type OutboxMessage struct {
	ID          string
	Subject     string
	Payload     []byte
	OccurredAt  time.Time
	PublishedAt *time.Time
}

type OutboxRepository interface {
	Add(ctx context.Context, msg OutboxMessage) error
}

type OutboxFactory interface {
	WithTx(tx uow.Tx) OutboxFactory
	Outbox() OutboxRepository
}
