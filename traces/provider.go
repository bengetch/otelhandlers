package traces

import (
	"context"
	"errors"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func GetExporter(exporterType string) (sdktrace.SpanExporter, error) {

	// send trace output to stdout
	if exporterType == "stdout" {
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	// send trace output to collector instance
	if exporterType == "otel" {

		// get endpoint for collector service from environment
		collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if collectorEndpoint == "" {
			return nil, errors.New(
				"failed to configure TracerProvider: exporterType set to `otel` " +
					"but environment variable OTEL_EXPORTER_OTLP_ENDPOINT is empty",
			)
		}
		return otlptrace.New(context.Background(), otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorEndpoint),
		))
	}

	// silence trace output
	return &NoOpSpanExporter{}, nil
}

func GetProvider(exporterType string, serviceName string) (*sdktrace.TracerProvider, error) {
	/*
		resolve which trace provider to use from the value of exporterType
	*/

	traceExporter, err := GetExporter(exporterType)
	if err != nil {
		return nil, err
	}

	// instantiate tracer provider with exporter defined above
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	return tracerProvider, nil
}
