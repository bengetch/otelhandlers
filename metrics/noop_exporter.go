package metrics

import (
	"context"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type NoOpMetricExporter struct{}

func (e *NoOpMetricExporter) Temporality(sdkmetric.InstrumentKind) metricdata.Temporality { return 0 }

func (e *NoOpMetricExporter) Aggregation(sdkmetric.InstrumentKind) sdkmetric.Aggregation { return nil }

func (e *NoOpMetricExporter) Export(context.Context, *metricdata.ResourceMetrics) error { return nil }

func (e *NoOpMetricExporter) ForceFlush(context.Context) error { return nil }

func (e *NoOpMetricExporter) Shutdown(context.Context) error { return nil }
