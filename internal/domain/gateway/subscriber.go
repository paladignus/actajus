// Package gateway
package gateway

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type MessageHandler func(ctx context.Context, event event.IEvent) error

type Subscriber interface {
	Subscribe(ctx context.Context, eventName string, handler MessageHandler) error
	Unsubscribe(eventName string) error
	Close() error
}
