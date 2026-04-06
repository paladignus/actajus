// Package interceptor provides gRPC/Connect interceptors for observability.
package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingInterceptor adds distributed tracing to gRPC/Connect calls.
type TracingInterceptor struct {
	serviceName string
}

// NewTracingInterceptor creates a new tracing interceptor.
func NewTracingInterceptor(serviceName string) *TracingInterceptor {
	return &TracingInterceptor{serviceName: serviceName}
}

// WrapUnary adds tracing to unary RPC calls.
func (t *TracingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		tracer := otel.Tracer(t.serviceName)
		proc := req.Spec().Procedure

		ctx, span := tracer.Start(
			ctx,
			proc,
			trace.WithAttributes(
				attribute.String("rpc.system", "connect"),
				attribute.String("rpc.stream", string(req.Spec().StreamType)),
				attribute.String("rpc.method", proc),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		resp, err := next(ctx, req)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		return resp, err
	}
}

// WrapStreamingClient adds tracing to streaming client calls.
func (t *TracingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		tracer := otel.Tracer(t.serviceName)

		ctx, span := tracer.Start(
			ctx,
			s.Procedure,
			trace.WithAttributes(
				attribute.String("rpc.system", "connect"),
				attribute.String("rpc.stream", "client"),
				attribute.String("rpc.stream_type", string(s.StreamType)),
			),
			trace.WithSpanKind(trace.SpanKindClient),
		)

		conn := next(ctx, s)
		span.End()

		return conn
	}
}

// WrapStreamingHandler adds tracing to streaming handler calls.
func (t *TracingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		tracer := otel.Tracer(t.serviceName)
		proc := shc.Spec().Procedure

		ctx, span := tracer.Start(
			ctx,
			proc,
			trace.WithAttributes(
				attribute.String("rpc.system", "connect"),
				attribute.String("rpc.stream", "handler"),
				attribute.String("rpc.stream_type", string(shc.Spec().StreamType)),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		err := next(ctx, shc)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		return err
	}
}
