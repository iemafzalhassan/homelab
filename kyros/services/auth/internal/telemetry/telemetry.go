package telemetry

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
)

// InitTelemetry initializes OpenTelemetry tracing and returns a tracer provider.
func InitTelemetry(serviceName string, collectorURL string) (*sdktrace.TracerProvider, error) {
	if collectorURL == "" {
		// If no collector URL is provided, return a no-op tracer provider.
		return sdktrace.NewTracerProvider(), nil
	}

	// Create the OTLP exporter.
	exporter, err := otlptrace.New(context.Background(), otlptrace.WithEndpoint(collectorURL), otlptrace.WithInsecure())
	if err != nil {
		return nil, err
	}

	// Create the resource.
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create the trace provider.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global tracer provider.
	otel.SetTracerProvider(tp)

	return tp, nil
}

// NewPrometheusMetrics returns a new Prometheus metrics handler.
// Note: This is a placeholder. In a real implementation, you would use a Prometheus library to create metrics.
// For now, we return a nil handler because we are using the existing metrics from kyros/internal/metrics.
func NewPrometheusMetrics() interface {
	Handler() http.Handler
} {
	// We are not implementing Prometheus metrics here because we are using the existing one from kyros/internal/metrics.
	// Return nil to avoid breaking the build. The main.go will need to be adjusted to use the existing metrics.
	return nil
}