// Package tracing provides HTTP middleware for distributed tracing
package tracing

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// responseWriterWrapper wraps the HTTP ResponseWriter to capture status code
// This allows us to record the HTTP status code in the trace span
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Allows the wrapper to support additional interfaces if needed
func (rw *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Middleware creates an HTTP middleware that adds tracing to incoming requests
// This middleware:
// 1. Extracts trace context from incoming request headers
// 2. Creates a new span for the HTTP request
// 3. Adds HTTP-specific attributes to the span
// 4. Sets span status based on the response status
func Middleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get tracer for this service
			tracer := otel.Tracer(serviceName)

			// Extract trace context from request headers and create new span
			// This connects the incoming request to potentially existing trace context
			ctx, span := tracer.Start(
				r.Context(),
				r.URL.Path, // Use the URL path as the span name
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
					attribute.String("http.remote_addr", r.RemoteAddr),
				),
				trace.WithSpanKind(trace.SpanKindServer), // Mark this as a server span
			)
			defer span.End() // Ensure span is ended when the request is handled

			// Replace request context with span context
			// This ensures that any operations within this request will be associated with this span
			r = r.WithContext(ctx)

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default to 200
			}

			// Record start time for calculating request duration
			start := time.Now()

			// Call the next handler in the chain with wrapped response writer
			next.ServeHTTP(wrapped, r)

			// Calculate duration and add to span
			duration := time.Since(start)
			span.SetAttributes(
				attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
				attribute.Int("http.status_code", wrapped.statusCode),
			)

			// Set span status based on HTTP status code
			if wrapped.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
		})
	}
}