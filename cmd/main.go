package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	metricInstrument "go.opentelemetry.io/otel/metric/instrument"
	"net/http"
	"otlp-stack/config"
	"otlp-stack/internal/telemetry"
	"otlp-stack/pkg/log"
)

func main() {
	l := log.Init(log.WithDevelopment(true), log.WithLevel(0))
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		l.V(5).Info("failed start configs")
	}
	/// Observability

	instrument, err := telemetry.InitTelemetry(ctx, l, cfg)
	if err != nil {
		l.V(5).Error(err, err.Error())
	}
	if err != nil {
		l.V(5).Info("%s: %v", "Failed to initialize opentelemetry provider", err)
	}

	StartGin(instrument)

}

func StartGin(inst telemetry.Instrument) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(otelgin.Middleware("test-otlp"))

	router.GET("/", func(c *gin.Context) {
		obsCounter, err := inst.Metric().Int64Counter(
			"request.count",
			metricInstrument.WithDescription("counting requests"),
			metricInstrument.WithUnit("1"))
		if err != nil {
			otel.Handle(err)
		}

		// Start a new span
		ctx, span := inst.StartRootSpan(c.Request.Context(), "my-gin-server.handler")
		defer span.End()
		ctx.Done()
		// Set attributes on the span
		span.SetAttributes(attribute.String("http.method", c.Request.Method), attribute.String("http.path", c.Request.URL.Path))

		// Do some work
		attrs := []attribute.KeyValue{
			attribute.String("path", "/"),
			attribute.String("method", c.Request.Method),
		}
		obsCounter.Add(ctx, 1, attrs...)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})

	router.Run(":5000")
}
