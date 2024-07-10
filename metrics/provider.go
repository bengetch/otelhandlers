package metrics

import (
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func GetExporter(exporterType string) (metric.Exporter, error) {

	// send meter output to stdout
	if exporterType == "stdout" {
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}

	// send meter output to collector instance
	if exporterType == "otel" {

		// get endpoint for collector service from environment
		collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if collectorEndpoint == "" {
			return nil, errors.New(
				"failed to configure MeterProvider: exporterType set to `otel` " +
					"but environment variable OTEL_EXPORTER_OTLP_ENDPOINT is empty",
			)
		}

		return otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpoint(collectorEndpoint),
		)
	}

	// silence meter output
	return &NoOpMetricExporter{}, nil
}

func GetProvider(exporterType string, serviceName string) (*metric.MeterProvider, error) {
	/*
		resolve which meter provider to use from the value of exporterType
	*/

	metricExporter, err := GetExporter(exporterType)
	if err != nil {
		return nil, err
	}

	// instantiate meter provider with exporter defined above
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(metricExporter, metric.WithInterval(10*time.Second)),
		),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	return meterProvider, nil
}
