// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/shared/application/messaging"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type OutboxRepository interface {
	Add(ctx context.Context, tx uow.Tx, msg messaging.OutboxMessage) error
}
