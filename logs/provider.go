package logs

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func GetExporter(exporterType string) (log.Exporter, error) {

	// send log output to stdout
	if exporterType == "stdout" {
		return stdoutlog.New(stdoutlog.WithPrettyPrint())
	}

	// send log output to collector (or datadog agent)
	if exporterType == "otel" {

		collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if collectorEndpoint == "" {
			return nil, errors.New(
				"failed to configure LoggerProvider: exporterType set to `otel` " +
					"but environment variable OTEL_EXPORTER_OTLP_ENDPOINT is empty",
			)
		}

		return otlploghttp.New(
			context.Background(),
			otlploghttp.WithInsecure(),
			otlploghttp.WithEndpoint(collectorEndpoint),
		)
	}

	// do not send log output anywhere
	return NoOpLogExporter{}, nil
}

func GetProvider(exporterType string, serviceName string) (*log.LoggerProvider, error) {
	/*
		resolve which log provider to use from the value of exporterType
	*/

	logExporter, err := GetExporter(exporterType)
	if err != nil {
		return nil, err
	}
	logProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)
	return logProvider, nil
}
