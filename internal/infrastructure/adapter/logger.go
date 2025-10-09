// Package adapter
package adapter

import (
	"context"
	"log/slog"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(logger *slog.Logger) repository.Logger {
	return &SlogAdapter{logger: logger}
}

func (s *SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	s.logger.DebugContext(ctx, msg, args...)
}

func (s *SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	s.logger.InfoContext(ctx, msg, args...)
}

func (s *SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	s.logger.WarnContext(ctx, msg, args...)
}

func (s *SlogAdapter) Error(ctx context.Context, msg string, args ...any) {
	s.logger.ErrorContext(ctx, msg, args...)
}

func (s *SlogAdapter) With(args ...any) repository.Logger {
	return &SlogAdapter{
		logger: s.logger.With(args...),
	}
}

func (s *SlogAdapter) WithError(err error) repository.Logger {
	return &SlogAdapter{
		logger: s.logger.With("error", err),
	}
}
