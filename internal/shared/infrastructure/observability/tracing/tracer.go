// Package tracing provides distributed tracing functionality using OpenTelemetry.
package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config holds the tracing configuration options.
type Config struct {
	ServiceName string
	EndpointURL string // OTLP endpoint URL (e.g., "http://jaeger:4318/v1/traces")
	Insecure    bool   // Whether to use insecure connection (no TLS)
}

// InitTracer creates a new trace provider instance using OTLP exporter and registers it as global.
// This enables distributed tracing for the application.
func InitTracer(cfg Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create OTLP HTTP exporter.
	// This exporter sends trace data to an OTLP-compatible collector (Jaeger, Tempo, etc.)
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpointURL(cfg.EndpointURL),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information.
	// This identifies the service in tracing data.
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
	)

	// Create trace provider with batch span processor.
	// Batch processor improves performance by sending spans in batches.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(rs),
	)

	// Set global trace provider.
	// This allows any part of the application to use tracing.
	otel.SetTracerProvider(tp)

	// Set propagator to handle trace context propagation.
	// This ensures trace IDs are passed between services/components.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Register shutdown function.
	// This ensures all pending spans are sent before program exits.
	log.Println("Tracing initialized for service:", cfg.ServiceName)
	return nil
}

// ShutdownTracer flushes all remaining spans and closes the exporter.
func ShutdownTracer(ctx context.Context) error {
	// Get the global tracer provider and cast it to TracerProvider interface.
	// The Shutdown method is available on the SDK's TracerProvider.
	if sdkProvider, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		return sdkProvider.Shutdown(ctx)
	}

	// If not using the SDK provider, just return nil.
	// This can happen in test environments or when using a mock provider.
	return nil
}
