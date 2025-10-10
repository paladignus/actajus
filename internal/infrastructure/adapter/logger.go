// Package adapter
package adapter

import (
	"context"
	"log/slog"

	"github.com/paladignus/actajus/internal/domain/repository"
	appctx "github.com/paladignus/actajus/internal/infrastructure/context"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(logger *slog.Logger) repository.Logger {
	return &SlogAdapter{logger: logger}
}

// extractContextFields extrai campos do contexto automaticamente
func (s *SlogAdapter) extractContextFields(ctx context.Context) []any {
	var fields []any

	if traceID := appctx.GetTraceID(ctx); traceID != "" {
		fields = append(fields, "trace_id", traceID)
	}

	if userID := appctx.GetUserID(ctx); userID != "" {
		fields = append(fields, "user_id", userID)
	}

	return fields
}

func (s *SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	contextFields := s.extractContextFields(ctx)
	allArgs := append(contextFields, args...)
	s.logger.DebugContext(ctx, msg, allArgs...)
}

func (s *SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	contextFields := s.extractContextFields(ctx)
	allArgs := append(contextFields, args...)
	s.logger.InfoContext(ctx, msg, allArgs...)
}

func (s *SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	contextFields := s.extractContextFields(ctx)
	allArgs := append(contextFields, args...)
	s.logger.WarnContext(ctx, msg, allArgs...)
}

func (s *SlogAdapter) Error(ctx context.Context, msg string, args ...any) {
	contextFields := s.extractContextFields(ctx)
	allArgs := append(contextFields, args...)
	s.logger.ErrorContext(ctx, msg, allArgs...)
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
