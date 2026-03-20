// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *Logger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *Logger) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *Logger) Error(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *Logger) With(keysAndValues ...any) repository.Logger {
	return m
}

func (m *Logger) WithError(err error) repository.Logger {
	return m
}
