package otelpp

import (
	"compress/gzip"
	"context"
	"github.com/pkg/errors"
	egzip "google.golang.org/grpc/encoding/gzip"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
)

func grpcMetricProvider(ctx context.Context, cfg Config) (*sdkmetric.MeterProvider, error) {
	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	timeout := 10 * time.Second
	if cfg.ValidTimeout() {
		timeout = cfg.Timeout
	}

	conn, err := createGrpcConn(ctx, cfg.MetricEndpoint, cfg.Insecure, cfg.GRPCBlockConn, timeout)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	exp, err := otlpmetricgrpc.New(ctx,
		withOtlpMetricGRPCOptions(cfg, conn)...,
	)

	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return createMetricProvider(res, exp, cfg)
}

func withOtlpMetricGRPCOptions(cfg Config, conn *grpc.ClientConn) []otlpmetricgrpc.Option {
	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithGRPCConn(conn)}

	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(cfg.Headers))
	}
	if cfg.ValidTimeout() {
		opts = append(opts, otlpmetricgrpc.WithTimeout(cfg.Timeout))
	}
	if cfg.RetryConfig != (RetryConfig{}) {
		opts = append(opts, otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
			Enabled:         cfg.Enabled,
			InitialInterval: cfg.InitialInterval,
			MaxInterval:     cfg.MaxInterval,
			MaxElapsedTime:  cfg.MaxElapsedTime,
		}))
	}

	if cfg.UseGzipCompression {
		opts = append(opts, otlpmetricgrpc.WithCompressor(egzip.Name))
		_ = egzip.SetLevel(gzip.BestSpeed)
	}

	return opts
}
