// Package gateway
package gateway

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type Publisher interface {
	Publish(ctx context.Context, event event.Event) error
	PublishBatch(ctx context.Context, events []event.Event) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, eventName string, handler event.Handler) error
	Unsubscribe(eventName string) error
	Close() error
}
