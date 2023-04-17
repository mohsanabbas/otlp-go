package otelpp

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	"compress/gzip"
	egzip "google.golang.org/grpc/encoding/gzip"
)

/*
NewGRPCProvider creates and sets the global trace and metric provider
configured with an OTel Exporter that exports the collected data via gRPC.

The returned Tracing and Metric structure can be used to create or shutdown
the processor.
*/
func NewGRPCProvider(ctx context.Context, opts ...OptionProvider) (tracing Telemetry, metric Meter, err error) {

	cfg := buildConfig(opts...)

	if hasMissingConfigInfo(cfg) {
		return nil, nil, ErrMissingConfig
	}

	setErrorHandler(cfg)

	if cfg.traceEnable() {
		tracing, err = newGRPCTracerProvider(ctx, cfg)
		if err != nil {
			return nil, nil, errors.Wrap(err, err.Error())
		}
	}

	if cfg.metricEnable() {
		metric, err = newGRPCMetricProvider(ctx, cfg)
		if err != nil {
			return nil, nil, errors.Wrap(err, err.Error())
		}
	}

	return
}

/*
newGRPCMetricProvider creates and sets the global trace provider
configured with an OTel Exporter that exports the collected data via gRPC.

The returned Metric structure can be used to create or shutdown
the processor.
*/
func newGRPCMetricProvider(ctx context.Context, cfg Config) (*Metric, error) {
	mp, err := grpcMetricProvider(ctx, cfg)
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
newGRPCTracerProvider creates and sets the global trace provider
configured with an OTel Exporter that exports the collected spans via gRPC.

The returned Tracing structure can be used to create new spans or shutdown
the span processor.
*/
func newGRPCTracerProvider(ctx context.Context, cfg Config) (*Tracing, error) {
	tp, err := grpcTraceProvider(ctx, cfg)
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

func grpcTraceProvider(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	res, err := createResource(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	timeout := 10 * time.Second
	if cfg.ValidTimeout() {
		timeout = cfg.Timeout
	}

	conn, err := createGrpcConn(ctx, cfg.TraceEndpoint, cfg.Insecure, cfg.GRPCBlockConn, timeout)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	exp, err := otlptracegrpc.New(ctx,
		withOtlpGRPCOptions(cfg, conn)...,
	)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return createTracerProvider(exp, res, cfg)
}

func withOtlpGRPCOptions(cfg Config, conn *grpc.ClientConn) []otlptracegrpc.Option {
	opts := []otlptracegrpc.Option{otlptracegrpc.WithGRPCConn(conn)}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
	}
	if cfg.ValidTimeout() {
		opts = append(opts, otlptracegrpc.WithTimeout(cfg.Timeout))
	}
	if cfg.RetryConfig != (RetryConfig{}) {
		opts = append(opts, otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         cfg.Enabled,
			InitialInterval: cfg.InitialInterval,
			MaxInterval:     cfg.MaxInterval,
			MaxElapsedTime:  cfg.MaxElapsedTime,
		}))
	}

	if cfg.UseGzipCompression {
		opts = append(opts, otlptracegrpc.WithCompressor(egzip.Name))
		_ = egzip.SetLevel(gzip.BestSpeed)
	}

	return opts
}

func createGrpcConn(ctx context.Context, endpoint string, insec, block bool, timeout time.Duration) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if insec {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if block {
		opts = append(opts, grpc.WithBlock())
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		trimEndpoint(endpoint),
		opts...,
	)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return conn, nil
}
