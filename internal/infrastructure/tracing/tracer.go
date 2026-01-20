// Package tracing provides distributed tracing functionality using OpenTelemetry
package tracing

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracer creates a new trace provider instance and registers it as global
// This enables distributed tracing for the application
func InitTracer(serviceName, agentHost, agentPort string) error {
	// Create Jaeger exporter
	// This exporter sends trace data to a Jaeger agent/collector
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(agentHost),
			jaeger.WithAgentPort(agentPort),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	// This identifies the service in tracing data
	rs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	// Create trace provider with batch span processor
	// Batch processor improves performance by sending spans in batches
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(rs),
	)

	// Set global trace provider
	// This allows any part of the application to use tracing
	otel.SetTracerProvider(tp)

	// Set propagator to handle trace context propagation
	// This ensures trace IDs are passed between services/components
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Register shutdown function
	// This ensures all pending spans are sent before program exits
	log.Println("Tracing initialized for service:", serviceName)
	return nil
}

// ShutdownTracer flushes all remaining spans and closes the exporter
func ShutdownTracer(ctx context.Context) error {
	// Get the global tracer provider and cast it to TracerProvider interface
	// The Shutdown method is available on the SDK's TracerProvider
	if sdkProvider, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		return sdkProvider.Shutdown(ctx)
	}

	// If not using the SDK provider, just return nil
	// This can happen in test environments or when using a mock provider
	return nil
}