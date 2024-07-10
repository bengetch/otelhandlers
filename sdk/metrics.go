package sdk

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bengetch/otelhandlers/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
)

func SetupMetrics(exporterType string, serviceName string) *metric.MeterProvider {
	/*
		configure meter provider instance, which is responsible for exporting metrics to the backend
		indicated by exoorterType
	*/

	mp, mpErr := metrics.GetProvider(exporterType, serviceName)
	if mpErr != nil {
		slog.ErrorContext(
			context.Background(),
			fmt.Sprintf("Failed to get metric provider: %v\n", mpErr),
		)
	} else {
		otel.SetMeterProvider(mp)
	}

	return mp
}
