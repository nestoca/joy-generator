package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type TracerOptions struct {
	ServiceName    string
	Environment    string
	ServiceVersion string
	OTLPEndpoint   string
}

func noopE() error { return nil }

func SetupTracer(opts TracerOptions) (func() error, error) {
	if opts.OTLPEndpoint == "" {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return noopE, nil
	}

	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(opts.OTLPEndpoint),
	)

	exporter, err := otlptrace.New(context.Background(), otlpClient)
	if err != nil {
		return nil, err
	}

	r, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		"",
		semconv.ServiceName(opts.ServiceName),
		semconv.ServiceVersion(opts.ServiceVersion),
		attribute.String("environment", opts.Environment),
	))
	if err != nil {
		return nil, err
	}

	provider := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(r))

	otel.SetTracerProvider(provider)

	teardown := func() error {
		return provider.Shutdown(context.Background())
	}

	return teardown, nil
}

// StartTrace creates a span using the default tracer instance.
func StartTrace(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return otel.Tracer("").Start(ctx, spanName, opts...)
}
