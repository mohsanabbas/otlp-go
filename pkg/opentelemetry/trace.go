package otelpp

//go:generate mockgen -destination=./mocks/mock_telemetry.go -package=mocks github.com/PicPay/lib-go-otel/v2 Telemetry
//go:generate mockgen -destination=./mocks/mock_span.go -package=mocks github.com/PicPay/lib-go-otel/v2 Span

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// Telemetry defines methods to handle spans and the span processor.
type Telemetry interface {
	Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span)
	Shutdown(ctx context.Context) error
}

// Span defines a spans' methods.
type Span interface {
	trace.Span
}

// Tracing is the structure to be used for handling OTel traces.
type Tracing struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

/*
Start creates a span and a context.Context containing the newly-created span.

Any Span that is created MUST also be ended. This is the responsibility of the user.
Implementations of this API may leak memory or other resources if Spans are not ended.
*/
func (t *Tracing) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// Shutdown shuts down the span processors in the order they were registered.
func (t *Tracing) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

func createResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.AppEnv.String()),
		),
	)
}

func createTracerProvider(e sdktrace.SpanExporter, r *resource.Resource, cfg Config) (*sdktrace.TracerProvider, error) {
	var opts []sdktrace.BatchSpanProcessorOption

	if cfg.sendIntervalTrace != nil {
		opts = append(opts, sdktrace.WithBatchTimeout(*cfg.sendIntervalTrace))
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(e, opts...),
		sdktrace.WithResource(r),
	), nil
}
