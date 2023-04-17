package otelpp

import (
	"github.com/go-logr/logr"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
)

type OptionProvider func(c *Config)
type TraceOptionProvider func(c *Config)
type MetricOptionProvider func(c *Config)
type RetryOptionProvider func(c *Config)

// WithTraceEndpoint - endpoint to send traces
func WithTraceEndpoint(endpoint string) OptionProvider {
	return func(c *Config) {
		c.TraceEndpoint = endpoint
	}
}

// WithMetricEndpoint - endpoint to send metrics
func WithMetricEndpoint(endpoint string) OptionProvider {
	return func(c *Config) {
		c.MetricEndpoint = endpoint
	}
}

// WithServiceName - service name
func WithServiceName(serviceName string) OptionProvider {
	return func(c *Config) {
		c.ServiceName = serviceName
	}
}

// WithAppEnv - app environment level
func WithAppEnv(appEnv EnvLevel) OptionProvider {
	return func(c *Config) {
		c.AppEnv = appEnv
	}
}

// WithHeaders - header to send with trace and metric
func WithHeaders(headers map[string]string) OptionProvider {
	return func(c *Config) {
		c.Headers = headers
	}
}

// WithInsecure - connection insecure
func WithInsecure(insecure bool) OptionProvider {
	return func(c *Config) {
		c.Insecure = insecure
	}
}

// WithGRPCConnectionBlock - gRPC blocking connection
func WithGRPCConnectionBlock(block bool) OptionProvider {
	return func(c *Config) {
		c.GRPCBlockConn = block
	}
}

// WithTimeout - connection timeout
func WithTimeout(timeout time.Duration) OptionProvider {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithMetric - metric options
func WithMetric(opts ...MetricOptionProvider) OptionProvider {
	return func(c *Config) {
		for _, o := range opts {
			o(c)
		}
	}
}

// WithTrace - trace options
func WithTrace(opts ...TraceOptionProvider) OptionProvider {
	return func(c *Config) {
		for _, o := range opts {
			o(c)
		}
	}
}

// WithMetricReader - set the reader for metric, otherwise use the default sdkmetric.NewPeriodicReader
func WithMetricReader(r metric.Reader) MetricOptionProvider {
	return func(c *Config) {
		c.reader = r
	}
}

// WithMetricViews - set the views for metric, otherwise use the default otelpp.NewMetricHistogramBucketView
func WithMetricViews(v []metric.View) MetricOptionProvider {
	return func(c *Config) {
		c.views = v
	}
}

// WithSendIntervalMetric - set send interval to otel collector
func WithSendIntervalMetric(si time.Duration) MetricOptionProvider {
	return func(c *Config) {
		c.sendIntervalMetric = &si
	}
}

// WithSendIntervalTrace - set send interval to otel collector
func WithSendIntervalTrace(si time.Duration) TraceOptionProvider {
	return func(c *Config) {
		c.sendIntervalTrace = &si
	}
}

// WithRetryDefault - retry options with default values
// Recommended at go.opentelemetry.io/otel/exporters/otlp/internal/retry.DefaultConfig
func WithRetryDefault() OptionProvider {
	return WithRetry(
		WithRetryEnable(true),
		WithRetryInitialInterval(5*time.Second),
		WithRetryMaxInterval(30*time.Second),
		WithRetryMaxElapsedTime(time.Minute),
	)
}

// WithRetry - retry options
func WithRetry(opts ...RetryOptionProvider) OptionProvider {
	return func(c *Config) {
		for _, opt := range opts {
			opt(c)
		}
	}
}

// WithRetryEnable - enable retry
func WithRetryEnable(enable bool) RetryOptionProvider {
	return func(c *Config) {
		c.RetryConfig.Enabled = enable
	}
}

// WithRetryInitialInterval - initial interval for retry
func WithRetryInitialInterval(t time.Duration) RetryOptionProvider {
	return func(c *Config) {
		c.RetryConfig.InitialInterval = t
	}
}

// WithRetryMaxInterval - max interval for retry
func WithRetryMaxInterval(t time.Duration) RetryOptionProvider {
	return func(c *Config) {
		c.RetryConfig.MaxInterval = t
	}
}

// WithRetryMaxElapsedTime - max elapsed time for retry
func WithRetryMaxElapsedTime(t time.Duration) RetryOptionProvider {
	return func(c *Config) {
		c.RetryConfig.MaxElapsedTime = t
	}
}

// WithLogger - add PicPay lib logger
func WithLogger(l logr.Logger) OptionProvider {
	return func(c *Config) {
		c.Logger = l
	}
}

// WithGzipCompression - set to use gzip compression - with best speed compression
func WithGzipCompression(useGzipCompression bool) OptionProvider {
	return func(c *Config) {
		c.UseGzipCompression = useGzipCompression
	}
}
