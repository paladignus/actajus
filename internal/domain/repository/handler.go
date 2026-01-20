// Package repository
package repository

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/event"
)

type IHandler interface {
	Handle(ctx context.Context, event event.IEvent) error
	CanHandle(event event.IEvent) bool
}
