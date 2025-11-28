// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type Subscriber struct {
	EventName      string
	SubscribeError error
}

func (m *Subscriber) Subscribe(ctx context.Context, eventName string, handler repository.IHandler) error {
	return m.SubscribeError
}

func (m *Subscriber) Unsubscribe(eventName string) error {
	return nil
}

func (m *Subscriber) Close() error {
	return nil
}
