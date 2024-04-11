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

func GetTracerProvider(exporterType string, serviceName string) (*sdktrace.TracerProvider, error) {
	/*
		resolve which trace provider to use from the value of exporterType
	*/

	var traceExporter sdktrace.SpanExporter
	var err error

	if exporterType == "stdout" {
		// export traces to stdout of this service
		traceExporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	} else if exporterType == "otel" {
		// get endpoint for collector service from environment
		collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if collectorEndpoint == "" {
			return nil, errors.New(
				"failed to configure TracerProvider: exporterType set to `otel` " +
					"but environment variable OTEL_EXPORTER_OTLP_ENDPOINT is empty",
			)
		}
		// export traces to an otel collector
		traceExporter, err = otlptrace.New(context.Background(), otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(collectorEndpoint),
		))
	} else {
		// if no exporter type is indicated, no spans will be exported
		traceExporter = &NoOpSpanExporter{}
	}
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
