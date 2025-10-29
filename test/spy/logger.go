// Package spy
package spy

import (
	"context"

	"github.com/paladignus/actajus/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

type SpyLogger struct {
	mock.Mock
}

func (m *SpyLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *SpyLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *SpyLogger) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *SpyLogger) Error(ctx context.Context, msg string, keysAndValues ...any) {
	m.Called(append([]any{ctx, msg}, keysAndValues...)...)
}

func (m *SpyLogger) With(keysAndValues ...any) repository.Logger {
	return m
}

func (m *SpyLogger) WithError(err error) repository.Logger {
	return m
}
