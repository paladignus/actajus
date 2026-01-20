// Package gateway
package gateway

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type Subscriber interface {
	Subscribe(ctx context.Context, eventName string, handler repository.IHandler) error
	Unsubscribe(eventName string) error
	Close() error
}
