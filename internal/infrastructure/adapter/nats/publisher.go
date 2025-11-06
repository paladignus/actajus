// Package nats
package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/paladignus/actajus/internal/domain/event"
)

type Publisher struct {
	js      nats.JetStreamContext
	timeout time.Duration
	dlqSubj string
}

func NewPublisher(
	js nats.JetStreamContext,
	timeout time.Duration,
	dlqSubj string,
) *Publisher {
	return &Publisher{js: js, timeout: timeout, dlqSubj: dlqSubj}
}

func (p *Publisher) Publish(ctx context.Context, evt event.DomainEvent) error {
	payload, err := evt.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	ack, err := p.js.Publish(evt.Subject(), payload, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to '%s' publish event: %w", evt.Subject(), err)
	}
	if ack == nil || ack.Sequence == 0 {
		return fmt.Errorf("no ack from stream")
	}
	return nil
}
