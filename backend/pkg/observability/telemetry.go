package observability

import (
	"context"
	"fmt"

	"document-reader-chatbot/configs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// InitTelemetry initializes OpenTelemetry with Jaeger tracing and Prometheus metrics
func InitTelemetry(cfg configs.ObservabilityConfig) (*trace.TracerProvider, *metric.MeterProvider, error) {
	// Create resource
	res, err := newResource(cfg.ServiceName, cfg.ServiceVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize tracing
	tracerProvider, err := initTracing(cfg, res)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	// Initialize metrics
	meterProvider, err := initMetrics(res)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Set global providers
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	// Set text map propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider, meterProvider, nil
}

// newResource creates a new OpenTelemetry resource
func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
}

// initTracing initializes OpenTelemetry tracing with OTLP HTTP exporter
func initTracing(cfg configs.ObservabilityConfig, res *resource.Resource) (*trace.TracerProvider, error) {
	// Create OTLP HTTP exporter
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(cfg.JaegerEndpoint),
		otlptracehttp.WithInsecure(), // Use insecure connection for development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	return tp, nil
}

// initMetrics initializes OpenTelemetry metrics with Prometheus exporter
func initMetrics(res *resource.Resource) (*metric.MeterProvider, error) {
	// Create Prometheus exporter
	exp, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exp),
	)

	return mp, nil
}

// StartSpan is a helper function to start a new span
func StartSpan(ctx context.Context, tracerName, operationName string) (context.Context, func()) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, operationName)
	return ctx, func() { span.End() }
}

// RecordError records an error in the current span
func RecordError(ctx context.Context, err error) {
	span := oteltrace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
}
