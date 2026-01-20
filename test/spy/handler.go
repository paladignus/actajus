// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type Handler struct{}

func (m *Handler) Handle(ctx context.Context, event event.IEvent) error {
	return nil
}

func (m *Handler) CanHandle(event event.IEvent) bool {
	return true
}
