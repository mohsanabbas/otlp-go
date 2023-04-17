package otelpp

import (
	"context"
	"fmt"
)

// Shutdown - call shutdown of tracer and metric
func Shutdown(ctx context.Context, tracing Telemetry, metric Meter) error {

	var (
		errT error
		errM error
	)

	if tracing != nil {
		errT = tracing.Shutdown(ctx)
	}

	if metric != nil {
		errM = metric.Shutdown(ctx)
	}

	if errT != nil && errM != nil {
		return fmt.Errorf("trace: %e, metric: %e", errT, errM)
	} else if errT != nil {
		return errT
	} else if errM != nil {
		return errM
	}

	return nil
}
