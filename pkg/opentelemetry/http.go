package otelpp

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

/*
NewHTTPProvider creates and sets the global trace and metric provider
configured with an OTel Exporter that exports the collected data via HTTP.

The returned Tracing and Metric structure can be used to create new spans or shutdown
the processor.
*/
func NewHTTPProvider(ctx context.Context, opts ...OptionProvider) (Telemetry, Meter, error) {
	var (
		tracing *Tracing
		metric  *Metric
		err     error
	)

	cfg := buildConfig(opts...)

	if hasMissingConfigInfo(cfg) {
		return nil, nil, ErrMissingConfig
	}

	setErrorHandler(cfg)

	if cfg.traceEnable() {
		tracing, err = newHTTPTracerProvider(ctx, cfg)
		if err != nil {
			return nil, nil, errors.Wrap(err, err.Error())
		}
	}

	if cfg.metricEnable() {
		metric, err = newHTTPMetricProvider(ctx, cfg)
		if err != nil {
			return nil, nil, errors.Wrap(err, err.Error())
		}
	}

	return tracing, metric, nil
}

func newHTTPMetricProvider(ctx context.Context, cfg Config) (*Metric, error) {
	mp, err := httpMetricExporter(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	global.SetMeterProvider(mp)
	meter := mp.Meter(cfg.ServiceName)

	return &Metric{
		provider: mp,
		meter:    meter,
	}, nil
}

/*
newHTTPTracerProvider creates and sets the global trace provider
configured with an OTel Exporter that exports the collected spans via HTTP.

The returned Tracing structure can be used to create new spans or shutdown
the span processor.
*/
func newHTTPTracerProvider(ctx context.Context, cfg Config) (*Tracing, error) {
	tp, err := httpTraceProvider(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(cfg.ServiceName)

	return &Tracing{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func httpTraceProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	exp, err := otlptracehttp.New(ctx,
		withOtlpTraceHTTPOptions(cfg)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return createTracerProvider(exp, res, cfg)
}

func withOtlpTraceHTTPOptions(cfg Config) []otlptracehttp.Option {
	endpoint := trimEndpoint(cfg.TraceEndpoint)
	opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}

	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(cfg.Headers))
	}
	if cfg.ValidTimeout() {
		opts = append(opts, otlptracehttp.WithTimeout(cfg.Timeout))
	}
	if cfg.RetryConfig != (RetryConfig{}) {
		opts = append(opts, otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         cfg.Enabled,
			InitialInterval: cfg.InitialInterval,
			MaxInterval:     cfg.MaxInterval,
			MaxElapsedTime:  cfg.MaxElapsedTime,
		}))
	}

	if cfg.UseGzipCompression {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	return opts
}
