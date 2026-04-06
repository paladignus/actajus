// Package outbox
package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paladignus/actajus/internal/shared/application/messaging"
)

type Dispatcher struct {
	pool      *pgxpool.Pool
	publisher messaging.Publisher
	limit     int
}

func NewDispatcher(
	pool *pgxpool.Pool,
	publisher messaging.Publisher,
	limit int,
) *Dispatcher {
	if limit <= 0 {
		limit = 100
	}
	return &Dispatcher{
		pool:      pool,
		publisher: publisher,
		limit:     limit,
	}
}

type pendingMessage struct {
	ID      string
	Subject string
	Payload []byte
}

func (d *Dispatcher) DispatchPending(ctx context.Context) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `
		SELECT idoutbox_messages, subject, payload
		FROM outbox_messages
		WHERE published_at IS NULL
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`, d.limit)
	if err != nil {
		return fmt.Errorf("query pending outbox: %w", err)
	}
	defer rows.Close()
	var items []pendingMessage
	for rows.Next() {
		var item pendingMessage
		if err := rows.Scan(&item.ID, &item.Subject, &item.Payload); err != nil {
			return fmt.Errorf("scan pending outbox: %w", err)
		}
		items = append(items, item)
	}
	for _, item := range items {
		if err := d.publisher.Publish(ctx, item.Subject, item.ID, item.Payload); err != nil {
			_, _ = tx.Exec(ctx, `
				UPDATE outbox_messages
				SET error_message = $2,
				    last_attempt_at = NOW(),
				    attempts = attempts + 1
				WHERE idoutbox_messages = $1
			`, item.ID, err.Error())
			return fmt.Errorf("publish outbox message %s: %w", item.ID, err)
		}
		_, err := tx.Exec(ctx, `
			UPDATE outbox_messages
			SET published_at = NOW(),
			    error_message = NULL,
			    last_attempt_at = NOW(),
			    attempts = attempts + 1
			WHERE idoutbox_messages = $1
		`, item.ID)
		if err != nil {
			return fmt.Errorf("mark outbox message published %s: %w", item.ID, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit outbox dispatch: %w", err)
	}
	return nil
}

func (d *Dispatcher) Run(ctx context.Context, interval time.Duration) error {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := d.DispatchPending(ctx); err != nil {
				// aqui você pluga teu logger
				fmt.Println("dispatch pending:", err)
			}
		}
	}
}
