// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type SpyPublisher struct {
	event event.IEvent
	Err   error
}

func (s *SpyPublisher) Publish(ctx context.Context, event event.IEvent) error {
	s.event = event
	return s.Err
}

func (s *SpyPublisher) PublishBatch(ctx context.Context, events []event.IEvent) error {
	return s.Err
}

func (s *SpyPublisher) Close() error {
	return s.Err
}
