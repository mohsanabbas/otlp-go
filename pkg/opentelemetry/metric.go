package otelpp

//go:generate mockgen -destination=./mocks/mock_meter.go -package=mocks github.com/PicPay/lib-go-otel/v2 Meter
//go:generate mockgen -destination=./mocks/mock_async_int64.go -package=mocks github.com/PicPay/lib-go-otel/v2 Int64ObservableCounter
//go:generate mockgen -destination=./mocks/mock_async_float64.go -package=mocks github.com/PicPay/lib-go-otel/v2 Float64ObservableCounter
//go:generate mockgen -destination=./mocks/mock_sync_int64.go -package=mocks github.com/PicPay/lib-go-otel/v2 Int64Counter
//go:generate mockgen -destination=./mocks/mock_sync_float64.go -package=mocks github.com/PicPay/lib-go-otel/v2 Float64Counter
//go:generate mockgen -destination=./mocks/mock_float64_observable_up_down_counter.go -package=mocks github.com/PicPay/lib-go-otel/v2 Float64ObservableUpDownCounter
//go:generate mockgen -destination=./mocks/mock_float64_observable_gauge.go -package=mocks github.com/PicPay/lib-go-otel/v2 Float64ObservableGauge
//go:generate mockgen -destination=./mocks/mock_int64_observable_up_down_counter.go -package=mocks github.com/PicPay/lib-go-otel/v2 Int64ObservableUpDownCounter
//go:generate mockgen -destination=./mocks/mock_int64_observable_gauge.go -package=mocks github.com/PicPay/lib-go-otel/v2 Int64ObservableGauge
//go:generate mockgen -destination=./mocks/mock_float64_up_down_counter.go -package=mocks github.com/PicPay/lib-go-otel/v2 Float64UpDownCounter
//go:generate mockgen -destination=./mocks/mock_int64_up_down_counter.go -package=mocks github.com/PicPay/lib-go-otel/v2 Int64UpDownCounter

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
	"time"
)

// Compile-time check Metric implements Meter.
var _ Meter = (*Metric)(nil)

// Meter - wrap interface of go.opentelemetry.io/otel/metric.Meter
type Meter interface {
	RegisterCallback(f metric.Callback, instruments ...instrument.Asynchronous) (metric.Registration, error)
	Shutdown(ctx context.Context) error
	Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (Float64ObservableCounter, error)
	Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (Float64ObservableUpDownCounter, error)
	Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (Float64ObservableGauge, error)
	Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (Int64ObservableCounter, error)
	Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (Int64ObservableUpDownCounter, error)
	Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (Int64ObservableGauge, error)
	Float64Counter(name string, options ...instrument.Float64Option) (Float64Counter, error)
	Float64UpDownCounter(name string, options ...instrument.Float64Option) (Float64UpDownCounter, error)
	Float64Histogram(name string, options ...instrument.Float64Option) (Float64Histogram, error)
	Int64Counter(name string, options ...instrument.Int64Option) (Int64Counter, error)
	Int64UpDownCounter(name string, options ...instrument.Int64Option) (Int64UpDownCounter, error)
	Int64Histogram(name string, options ...instrument.Int64Option) (Int64Histogram, error)
}

type Float64ObservableCounter interface {
	instrument.Float64ObservableCounter
}

type Float64ObservableUpDownCounter interface {
	instrument.Float64ObservableUpDownCounter
}

type Float64ObservableGauge interface {
	instrument.Float64ObservableGauge
}

type Int64ObservableCounter interface {
	instrument.Int64ObservableCounter
}

type Int64ObservableUpDownCounter interface {
	instrument.Int64ObservableUpDownCounter
}

type Int64ObservableGauge interface {
	instrument.Int64ObservableGauge
}

type Float64Counter interface {
	instrument.Float64Counter
}

type Float64UpDownCounter interface {
	instrument.Float64UpDownCounter
}

type Float64Histogram interface {
	instrument.Float64Histogram
}

type Int64Counter interface {
	instrument.Int64Counter
}

type Int64UpDownCounter interface {
	instrument.Int64UpDownCounter
}

type Int64Histogram interface {
	instrument.Int64Histogram
}

// Metric is the structure to be used for handling OTel metrics.
type Metric struct {
	provider *sdkmetric.MeterProvider
	meter    metric.Meter
}

// Float64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing float64 measurements once per a measurement collection cycle.
func (m *Metric) Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (Float64ObservableCounter, error) {
	return m.meter.Float64ObservableCounter(name, options...)
}

// Float64ObservableUpDownCounter returns a new instrument identified by
// name and configured with options. The instrument is used to
// asynchronously record float64 measurements once per a measurement
// collection cycle.
func (m *Metric) Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (Float64ObservableUpDownCounter, error) {
	return m.meter.Float64ObservableUpDownCounter(name, options...)
}

// Float64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous float64 measurements once per a measurement collection
// cycle.
func (m *Metric) Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (Float64ObservableGauge, error) {
	return m.meter.Float64ObservableGauge(name, options...)
}

// Int64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing int64 measurements once per a measurement collection cycle.
func (m *Metric) Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (Int64ObservableCounter, error) {
	return m.meter.Int64ObservableCounter(name, options...)
}

// Int64ObservableUpDownCounter returns a new instrument identified by name
// and configured with options. The instrument is used to asynchronously
// record int64 measurements once per a measurement collection cycle.
func (m *Metric) Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (Int64ObservableUpDownCounter, error) {
	return m.meter.Int64ObservableUpDownCounter(name, options...)
}

// Int64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous int64 measurements once per a measurement collection
// cycle.
func (m *Metric) Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (Int64ObservableGauge, error) {
	return m.meter.Int64ObservableGauge(name, options...)
}

// Float64Counter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// increasing float64 measurements during a computational operation.
func (m *Metric) Float64Counter(name string, options ...instrument.Float64Option) (Float64Counter, error) {
	return m.meter.Float64Counter(name, options...)
}

// Float64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// float64 measurements during a computational operation.
func (m *Metric) Float64UpDownCounter(name string, options ...instrument.Float64Option) (Float64UpDownCounter, error) {
	return m.meter.Float64UpDownCounter(name, options...)
}

// Float64Histogram returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// the distribution of float64 measurements during a computational
// operation.
func (m *Metric) Float64Histogram(name string, options ...instrument.Float64Option) (Float64Histogram, error) {
	return m.meter.Float64Histogram(name, options...)
}

// Int64Counter returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record increasing
// int64 measurements during a computational operation.
func (m *Metric) Int64Counter(name string, options ...instrument.Int64Option) (Int64Counter, error) {
	return m.meter.Int64Counter(name, options...)
}

// Int64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// int64 measurements during a computational operation.
func (m *Metric) Int64UpDownCounter(name string, options ...instrument.Int64Option) (Int64UpDownCounter, error) {
	return m.meter.Int64UpDownCounter(name, options...)
}

// Int64Histogram returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// the distribution of int64 measurements during a computational operation.
func (m *Metric) Int64Histogram(name string, options ...instrument.Int64Option) (Int64Histogram, error) {
	return m.meter.Int64Histogram(name, options...)
}

// RegisterCallback captures the function that will be called during Collect.
func (m *Metric) RegisterCallback(f metric.Callback, instruments ...instrument.Asynchronous) (metric.Registration, error) {
	return m.meter.RegisterCallback(f, instruments...)
}

// Shutdown shuts down the span processors in the order they were registered.
func (m *Metric) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}

func createMetricProvider(r *resource.Resource, exp sdkmetric.Exporter, cfg Config) (*sdkmetric.MeterProvider, error) {

	reader := createMetricReader(exp, cfg)
	views := createMetricViews(cfg)

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(r),
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(views...),
	), nil
}

func createMetricReader(e sdkmetric.Exporter, cfg Config) sdkmetric.Reader {
	var opts []sdkmetric.PeriodicReaderOption

	if cfg.sendIntervalMetric != nil {
		opts = append(opts, sdkmetric.WithInterval(*cfg.sendIntervalMetric))
	}

	if cfg.ValidTimeout() {
		opts = append(opts, sdkmetric.WithTimeout(cfg.Timeout))
	}

	if cfg.MetricConfig.reader != nil {
		return cfg.MetricConfig.reader
	}

	return sdkmetric.NewPeriodicReader(e, opts...)
}

func createMetricViews(cfg Config) []sdkmetric.View {
	if len(cfg.MetricConfig.views) > 0 {
		return cfg.MetricConfig.views
	}

	return NewMetricHistogramBucketView()
}

func NewMetricHistogramBucketView() []sdkmetric.View {
	return []sdkmetric.View{
		sdkmetric.NewView(
			sdkmetric.Instrument{Kind: sdkmetric.InstrumentKindHistogram},
			sdkmetric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{
					float64(500 * time.Millisecond.Milliseconds()),
					float64(1 * time.Second.Milliseconds()),
					float64(10 * time.Second.Milliseconds()),
					float64(30 * time.Second.Milliseconds()),
					float64(1 * time.Minute.Milliseconds()),
				},
			}},
		),
	}
}
