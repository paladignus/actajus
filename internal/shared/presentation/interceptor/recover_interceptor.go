// Package interceptor
package interceptor

import (
	"context"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/domain/repository"
)

type RecoverInterceptor struct {
	logger repository.Logger
}

func NewRecoverInterceptor(logger repository.Logger) *RecoverInterceptor {
	return &RecoverInterceptor{logger}
}

func (r *RecoverInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				r.logger.Error(ctx,
					"panic recovered",
					"procedure", req.Spec().Procedure,
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				resp = nil
				err = connect.NewError(connect.CodeInternal, nil)
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
				r.logger.Error(ctx,
					"panic recovered",
					"procedure", shc.Spec().Procedure,
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				err = connect.NewError(connect.CodeInternal, nil)
			}
		}()
		return next(ctx, shc)
	}
}
