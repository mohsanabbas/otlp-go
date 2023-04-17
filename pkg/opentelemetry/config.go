package otelpp

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"strings"
	"time"
)

type EnvLevel int

const (
	INVALID EnvLevel = iota
	DEV
	QA
	PROD
)

var toString = map[EnvLevel]string{
	DEV:  "dev",
	QA:   "qa",
	PROD: "prod",
}

var (
	ErrMissingConfig       = errors.New("missing required fields: AppEnv, Endpoint (Metric and/or Trace), ServiceName")
	ErrMissingJaegerConfig = errors.New("missing required fields: AppEnv, TraceEndpoint, ServiceName")
)

// Config struct defines the required fields to create tracer providers.
type Config struct {
	AppEnv         EnvLevel
	TraceEndpoint  string
	MetricEndpoint string
	ServiceName    string
	Logger         logr.Logger
	JaegerConfig
	OtlpConfig
}

func (c *Config) traceEnable() bool {
	return c.TraceEndpoint != ""
}

func (c *Config) metricEnable() bool {
	return c.MetricEndpoint != ""
}

// JaegerConfig contains specific field for the Jaeger tracer provider
type JaegerConfig struct {
}

// OtlpConfig contains specific field for gRPC and HTTP tracer providers
type OtlpConfig struct {
	Headers            map[string]string
	Insecure           bool
	Timeout            time.Duration
	GRPCBlockConn      bool
	UseGzipCompression bool
	RetryConfig
	MetricConfig
	TraceConfig
}

// MetricConfig - configuration for metric
// View is an override to the default behavior of the SDK. It defines how data
// should be collected for certain instruments. use default otelpp.createMetricHistogramBucketView()
// SendIntervalMetric - default value 60s, defined at sdk metric.defaultInterval
type MetricConfig struct {
	sendIntervalMetric *time.Duration
	reader             metric.Reader
	views              []metric.View
}

// TraceConfig - configuration for trace
// SendIntervalTrace - default value 5s, defined at sdk trace.DefaultScheduleDelay
type TraceConfig struct {
	sendIntervalTrace *time.Duration
}

func (c *OtlpConfig) ValidTimeout() bool {
	return c.Timeout.String() != "0s"
}

// RetryConfig contains the fields used by the internal package retry.Config,
// defined by go.opentelemetry.io/otel/exporters/otlp/internal/retry
type RetryConfig struct {
	Enabled         bool
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
}

func hasMissingConfigInfo(cfg Config) bool {
	return (cfg.AppEnv < 1 || cfg.AppEnv > 3) ||
		cfg.ServiceName == "" ||
		(!cfg.traceEnable() && !cfg.metricEnable())
}

// String used to translate an EnvLevel to string
func (l EnvLevel) String() string {
	if value, ok := toString[l]; ok {
		return value
	}
	return fmt.Sprintf("UNKNOWN[%d]", l)
}

func EnvLevelFromString(envLevel string) (EnvLevel, error) {
	envLevel = strings.ToLower(envLevel)
	for el, str := range toString {
		if str == envLevel && el != INVALID {
			return el, nil
		}
	}
	return INVALID, fmt.Errorf("invalid env level")
}

// trimEndpoint removes the 'http://' portion of the endpoint
func trimEndpoint(endpoint string) string {
	return strings.ReplaceAll(endpoint, "http://", "")
}

func buildConfig(opts ...OptionProvider) Config {
	cfg := Config{}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

func setErrorHandler(cfg Config) {
	l := cfg.Logger

	//l = log.Init(log.WithDevelopment(true), log.WithLevel(0))

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		l.Error(err, "OTel SDK error handler")
	}))
}
