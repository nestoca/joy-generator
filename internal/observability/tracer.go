package observability

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
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
		otlptracegrpc.WithTimeout(90*time.Second),
		// Configure gRPC connection retry for network-level failures (e.g., "connection refused")
		otlptracegrpc.WithDialOption(
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					// this delay/multiplier config gives retries at roughly [3, 9, 21, 45, 93] seconds
					BaseDelay:  3 * time.Second,
					Multiplier: 2.0,
					Jitter:     0.2,
					MaxDelay:   120 * time.Second,
				},
				MinConnectTimeout: 5 * time.Second,
			}),
		),
		// Configure application-level retry for retryable errors (e.g., rate limits, temporary server errors)
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 5 * time.Second,
			MaxInterval:     15 * time.Second,
			MaxElapsedTime:  90 * time.Second,
		}),
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
