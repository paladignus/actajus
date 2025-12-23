package repository

import (
	"context"
	"testing"
)

type LoggerSpy struct {
	calls []string
}

func (l *LoggerSpy) Debug(ctx context.Context, msg string, args ...any) {
	l.calls = append(l.calls, "Debug")
}

func (l *LoggerSpy) Info(ctx context.Context, msg string, args ...any) {
	l.calls = append(l.calls, "Info")
}

func (l *LoggerSpy) Warn(ctx context.Context, msg string, args ...any) {
	l.calls = append(l.calls, "Warn")
}

func (l *LoggerSpy) Error(ctx context.Context, msg string, args ...any) {
	l.calls = append(l.calls, "Error")
}

func (l *LoggerSpy) With(args ...any) Logger {
	l.calls = append(l.calls, "With")
	return l
}

func (l *LoggerSpy) WithError(err error) Logger {
	l.calls = append(l.calls, "WithError")
	return l
}

func TestLoggerInterface(t *testing.T) {
	ctx := context.Background()
	sut := &LoggerSpy{}
	sut.Debug(ctx, "debug message")
	sut.Info(ctx, "info message")
	sut.Warn(ctx, "warn message")
	sut.Error(ctx, "error message")
	_ = sut.With("key", "value")
	_ = sut.WithError(nil)
	expectedCalls := []string{"Debug", "Info", "Warn", "Error", "With", "WithError"}
	for i, expectedCall := range expectedCalls {
		if i < len(sut.calls) && sut.calls[i] != expectedCall {
			t.Errorf("Expected call %s, got %s", expectedCall, sut.calls[i])
		}
	}
}
