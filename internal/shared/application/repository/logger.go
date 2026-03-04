// Package repository
package repository

import "context"

type Logger interface {
	Debug(context.Context, string, ...any)
	Info(context.Context, string, ...any)
	Warn(context.Context, string, ...any)
	Error(context.Context, string, ...any)
	With(...any) Logger
	WithError(error) Logger
}
