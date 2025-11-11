// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type SpyPublisher struct {
	event event.Event
	Err   error
}

func (s *SpyPublisher) Publish(ctx context.Context, event event.Event) error {
	s.event = event
	return s.Err
}

func (s *SpyPublisher) PublishBatch(ctx context.Context, events []event.Event) error {
	return s.Err
}

func (s *SpyPublisher) Close() error {
	return s.Err
}
