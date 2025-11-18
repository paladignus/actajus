// Package event
package event

import "context"

type Handler interface {
	Handle(ctx context.Context, event IEvent) error
	CanHandle(event IEvent) bool
}
