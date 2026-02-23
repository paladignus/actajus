// Package interceptor
package interceptor

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
)

type LoggingInterceptor struct {
	logger repository.Logger
}

func NewLoggingInterceptor(logger repository.Logger) *LoggingInterceptor {
	return &LoggingInterceptor{logger: logger}
}

func (l *LoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		proc := req.Spec().Procedure
		method := req.HTTPMethod()
		l.logger.Info(ctx, "→ "+proc, "method", method)
		resp, err := next(ctx, req)
		dur := time.Since(start)
		if err != nil {
			l.logger.Error(ctx, "← "+proc,
				"method", method,
				"status", "error",
				"took", dur.String(),
				"error", err,
			)
			return nil, err
		}
		l.logger.Info(ctx, "← "+proc,
			"method", method,
			"status", "ok",
			"took", dur.String(),
		)
		return resp, nil
	}
}

func (l *LoggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (l *LoggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		start := time.Now()
		proc := shc.Spec().Procedure
		l.logger.Info(ctx, "→ "+proc, "stream", true)
		err := next(ctx, shc)
		dur := time.Since(start)
		if err != nil {
			l.logger.Error(ctx, "← "+proc,
				"stream", true,
				"status", "error",
				"took", dur.String(),
				"error", err,
			)
			return err
		}
		l.logger.Info(ctx, "← "+proc,
			"stream", true,
			"status", "ok",
			"took", dur.String(),
		)
		return nil
	}
}
