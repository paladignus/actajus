// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type SpyPublisher struct {
	event event.DomainEvent
	Err   error
}

func (s *SpyPublisher) Publish(ctx context.Context, event event.DomainEvent) error {
	s.event = event
	return s.Err
}
