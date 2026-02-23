// Package interceptor
package interceptor

import (
	"context"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
)

type RecoverInterceptor struct {
	logger *slog.Logger
}

func NewRecoverInterceptor(logger *slog.Logger) *RecoverInterceptor {
	return &RecoverInterceptor{logger: logger}
}

func (r *RecoverInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				if r.logger != nil {
					r.logger.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("procedure", req.Spec().Procedure),
						slog.String("stack", string(debug.Stack())),
					)
				}
				err = connect.NewError(connect.CodeInternal, nil)
				resp = nil
			}
		}()
		return next(ctx, req)
	}
}

func (r *RecoverInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (r *RecoverInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				if r.logger != nil {
					r.logger.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("procedure", shc.Spec().Procedure),
						slog.String("stack", string(debug.Stack())),
					)
				}
				err = connect.NewError(connect.CodeInternal, nil)
			}
		}()
		return next(ctx, shc)
	}
}
