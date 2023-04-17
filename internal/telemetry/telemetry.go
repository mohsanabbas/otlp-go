package telemetry

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/trace"
	"otlp-stack/config"
	otelpp "otlp-stack/pkg/opentelemetry"
	"time"
)

const (
	ErrStartInstrumentTimeout = "timeout to start instrumentation"
	ErrStartInstrument        = "failed to start instrument"
	ErrInvalidAppEnv          = "invalid application environment"
	otlTimeout                = time.Second * 5
)

type Instrument interface {
	StartRootSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	Metric() otelpp.Meter
}
type instrument struct {
	config *config.Config
	log    logr.Logger
	trace  otelpp.Telemetry
	metric otelpp.Meter
}

func (i instrument) StartRootSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := i.trace.Start(
		ctx,
		name,
		opts...,
	)
	return trace.ContextWithSpan(ctx, span), span
}

func (i instrument) Metric() otelpp.Meter {
	return i.metric
}

func InitTelemetry(ctx context.Context, l logr.Logger, cfg *config.Config) (Instrument, error) {
	appEnv, err := otelpp.EnvLevelFromString(cfg.AppStage)
	if err != nil {
		//errMsg := fmt.Sprintf("%s. Current: \"%v\"", ErrInvalidAppEnv, cfg.AppStage)

		return nil, err
	}

	instrumentation := instrument{
		config: cfg,
		log:    l,
	}

	ctxTimeout, cancelCtx := context.WithTimeout(ctx, otlTimeout)
	defer func() {
		cancelCtx()
		if instrumentation.trace == nil || instrumentation.metric == nil {
			instrumentation.log.Error(errors.New(ErrStartInstrumentTimeout), ErrStartInstrumentTimeout)
		}
	}()

	instrumentation.trace, instrumentation.metric, err = otelpp.NewGRPCProvider(ctxTimeout,
		otelpp.WithAppEnv(appEnv),
		otelpp.WithTraceEndpoint(cfg.TraceHost),
		otelpp.WithTrace(),
		otelpp.WithMetricEndpoint(cfg.MetricHost),
		otelpp.WithMetric(),
		otelpp.WithServiceName(cfg.ServiceName),
		otelpp.WithInsecure(true),
		otelpp.WithRetryDefault(),
		otelpp.WithTimeout(otlTimeout),
		otelpp.WithLogger(l),
		otelpp.WithGzipCompression(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	if err = metricProvider(); err != nil {
		return nil, errors.Wrap(err, err.Error())
	}

	return &instrumentation, nil
}

func metricProvider() error {
	if err := host.Start(); err != nil {
		return errors.Wrap(err, err.Error())
	}

	if err := runtime.Start(); err != nil {
		return errors.Wrap(err, err.Error())
	}

	return nil
}
