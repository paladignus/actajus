// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type Publisher struct {
	Event event.IEvent
	Err   error
}

func (s *Publisher) Publish(ctx context.Context, event event.IEvent) error {
	s.Event = event
	return s.Err
}

func (s *Publisher) PublishBatch(ctx context.Context, events []event.IEvent) error {
	return s.Err
}

func (s *Publisher) Close() error {
	return s.Err
}

// type Publisher struct {
// 	PublishError error
// 	PublishEvent event.IEvent
// }
//
// func (p *Publisher) Publish(ctx context.Context, event event.IEvent) error {
// 	p.PublishEvent = event
// 	return p.PublishError
// }
//
// func (p *Publisher) PublishBatch(ctx context.Context, events []event.IEvent) error {
// 	return nil
// }
//
// func (p *Publisher) Close() error {
// 	return nil
// }
