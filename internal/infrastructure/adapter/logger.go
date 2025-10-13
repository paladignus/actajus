// Package adapter
package adapter

import (
	"context"
	"log/slog"
	"os"

	"github.com/paladignus/actajus/internal/domain/repository"
)

type Level int

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type slogAdapter struct {
	logger *slog.Logger
}

type Config struct {
	Level     Level
	AddSource bool
	Format    string // "json" ou "text"
	Output    string // "stdout", "stderr" ou caminho do arquivo
}

func NewLogger(config Config) repository.Logger {
	var level slog.Level
	switch level {
	case slog.LevelDebug:
		level = slog.LevelDebug
	case slog.LevelInfo:
		level = slog.LevelInfo
	case slog.LevelWarn:
		level = slog.LevelWarn
	case slog.LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	var output = os.Stdout

	switch config.Output {
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stdout
	}

	switch config.Format {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return &slogAdapter{
		logger: slog.New(handler),
	}
}

func NewDefaultLogger() repository.Logger {
	return NewLogger(Config{
		Level:     LevelInfo,
		AddSource: false,
		Format:    "json",
		Output:    "stdout",
	})
}

func (s *slogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	s.logger.DebugContext(ctx, msg, args...)
}

func (s *slogAdapter) Info(ctx context.Context, msg string, args ...any) {
	s.logger.InfoContext(ctx, msg, args...)
}

func (s *slogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	s.logger.WarnContext(ctx, msg, args...)
}

func (s *slogAdapter) Error(ctx context.Context, msg string, args ...any) {
	s.logger.ErrorContext(ctx, msg, args...)
}

func (s *slogAdapter) With(args ...any) repository.Logger {
	return &slogAdapter{
		logger: s.logger.With(args...),
	}
}

func (s *slogAdapter) WithError(err error) repository.Logger {
	return &slogAdapter{
		logger: s.logger.With("error", err),
	}
}
