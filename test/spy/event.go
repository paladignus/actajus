// Package spy
package spy

import (
	"time"
)

type Event struct{}

func (m Event) EventName() string {
	return "mock.event"
}

func (m Event) OccurredAt() time.Time {
	return time.Now()
}

func (m Event) EventVersion() string {
	return "1.0"
}

func (m Event) AggregateID() string {
	return "mock-aggregate-id"
}

// type Publisher struct{}
//
// func (m *Publisher) Publish(ctx context.Context, event event.IEvent) error {
// 	return nil
// }
//
// func (m *Publisher) PublishBatch(ctx context.Context, events []event.IEvent) error {
// 	return nil
// }
//
// func (m *Publisher) Close() error {
// 	return nil
// }
