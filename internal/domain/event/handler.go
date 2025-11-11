// Package event
package event

import "context"

type Handler interface {
	Handle(ctx context.Context, event Event) error
	CanHandle(event Event) bool
}

type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

func (f HandlerFunc) CanHandle(event Event) bool {
	return true
}
