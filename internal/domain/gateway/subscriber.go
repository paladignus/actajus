// Package gateway
package gateway

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type Subscriber interface {
	Subscribe(ctx context.Context, eventName string, handler event.Handler) error
	Unsubscribe(eventName string) error
	Close() error
}
