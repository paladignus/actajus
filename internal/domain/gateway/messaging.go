// Package gateway
package gateway

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type Publisher interface {
	Publish(ctx context.Context, event event.DomainEvent) error
}
