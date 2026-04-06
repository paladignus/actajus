// Package interceptor provides gRPC/Connect interceptors for observability.
package interceptor

import (
	"context"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/infrastructure/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsInterceptor collects Prometheus metrics for gRPC/Connect calls.
type MetricsInterceptor struct{}

// NewMetricsInterceptor creates a new metrics interceptor.
func NewMetricsInterceptor() *MetricsInterceptor {
	return &MetricsInterceptor{}
}

// WrapUnary collects metrics for unary RPC calls.
func (m *MetricsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		start := time.Now()
		proc := req.Spec().Procedure

		resp, err := next(ctx, req)

		duration := time.Since(start)
		statusCode := "ok"
		if err != nil {
			statusCode = "error"
		}

		// Record gRPC request count
		metrics.RequestCount.With(prometheus.Labels{
			"method": req.HTTPMethod(),
			"path":   proc,
			"status": statusCode,
		}).Inc()

		// Record gRPC request duration
		metrics.RequestDuration.With(prometheus.Labels{
			"method": req.HTTPMethod(),
			"path":   proc,
			"status": statusCode,
		}).Observe(duration.Seconds())

		return resp, err
	}
}

// WrapStreamingClient collects metrics for streaming client calls.
func (m *MetricsInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		start := time.Now()

		conn := next(ctx, s)

		duration := time.Since(start)

		metrics.RequestCount.With(prometheus.Labels{
			"method": "STREAM",
			"path":   s.Procedure,
			"status": "ok",
		}).Inc()

		metrics.RequestDuration.With(prometheus.Labels{
			"method": "STREAM",
			"path":   s.Procedure,
			"status": "ok",
		}).Observe(duration.Seconds())

		return conn
	}
}

// WrapStreamingHandler collects metrics for streaming handler calls.
func (m *MetricsInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		start := time.Now()
		proc := shc.Spec().Procedure

		err := next(ctx, shc)

		duration := time.Since(start)
		statusCode := "ok"
		if err != nil {
			statusCode = "error"
		}

		metrics.RequestCount.With(prometheus.Labels{
			"method": "STREAM",
			"path":   proc,
			"status": statusCode,
		}).Inc()

		metrics.RequestDuration.With(prometheus.Labels{
			"method": "STREAM",
			"path":   proc,
			"status": statusCode,
		}).Observe(duration.Seconds())

		return err
	}
}

// GRPCCodeToString converts a Connect error code to a string.
func GRPCCodeToString(code connect.Code) string {
	return strconv.Itoa(int(code))
}
