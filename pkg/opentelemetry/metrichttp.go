package otelpp

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func httpMetricExporter(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, error) {
	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	exp, err := otlpmetrichttp.New(ctx, withOtlpMetricHTTPOptions(cfg)...)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return createMetricProvider(res, exp, cfg)
}

func withOtlpMetricHTTPOptions(cfg Config) []otlpmetrichttp.Option {
	endpoint := trimEndpoint(cfg.MetricEndpoint)
	opts := []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint(endpoint)}

	if cfg.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(cfg.Headers))
	}
	if cfg.ValidTimeout() {
		opts = append(opts, otlpmetrichttp.WithTimeout(cfg.Timeout))
	}
	if cfg.RetryConfig != (RetryConfig{}) {
		opts = append(opts, otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{
			Enabled:         cfg.Enabled,
			InitialInterval: cfg.InitialInterval,
			MaxInterval:     cfg.MaxInterval,
			MaxElapsedTime:  cfg.MaxElapsedTime,
		}))
	}

	if cfg.UseGzipCompression {
		opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	}

	return opts
}
