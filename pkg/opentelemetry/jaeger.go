package otelpp

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

/*
NewJaegerTracerProvider creates and sets the global trace provider
configured with an OTel Exporter that exports the collected spans to Jaeger.

The returned Tracing structure can be used to create new spans or shutdown
the span processor.
*/
func NewJaegerTracerProvider(ctx context.Context, cfg Config) (*Tracing, error) {
	if !cfg.traceEnable() || hasMissingConfigInfo(cfg) {
		return nil, ErrMissingJaegerConfig
	}

	setErrorHandler(cfg)

	tp, err := jaegerTraceProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(cfg.ServiceName)

	return &Tracing{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func jaegerTraceProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.TraceEndpoint),
		),
	)
	if err != nil {
		return nil, err
	}

	return createTracerProvider(exp, res, cfg)
}
