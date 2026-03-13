// Package postgres
package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
	"github.com/paladignus/actajus/internal/shared/application/uow"
)

type OutboxFactory struct {
	exec Executor
}

func NewOutboxFactory(exec Executor) *OutboxFactory {
	return &OutboxFactory{exec: exec}
}

func (f *OutboxFactory) WithTx(tx uow.Tx) messaging.OutboxFactory {
	pgxTx, ok := tx.(pgx.Tx)
	if !ok {
		panic(fmt.Sprintf("invalid tx type: %T", tx))
	}
	return &OutboxFactory{exec: pgxTx}
}

func (f *OutboxFactory) Outbox() messaging.OutboxRepository {
	return &OutboxRepository{exec: f.exec}
}
