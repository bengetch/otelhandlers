package logs

import (
	"context"
	"errors"
	"os"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func GetLogProvider(exporterType string, serviceName string) (*sdk.LoggerProvider, error) {
	/*
		resolve which log provider to use from the value of exporterType
	*/

	var logExporter sdk.LogRecordExporter
	var err error

	if exporterType == "stdout" {
		// export logs to stdout of this service
		logExporter, err = stdoutlogs.NewExporter()

	} else if exporterType == "otel" {
		// get endpoint for collector service from environment
		collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if collectorEndpoint == "" {
			return nil, errors.New(
				"failed to configure LoggerProvider: exporterType set to `otel` " +
					"but environment variable OTEL_EXPORTER_OTLP_ENDPOINT is empty",
			)
		}
		// export logs to an otel collector
		logExporter, err = otlplogs.NewExporter(
			context.Background(),
			otlplogs.WithClient(otlplogsgrpc.NewClient(
				otlplogsgrpc.WithInsecure(),
				otlplogsgrpc.WithEndpoint(collectorEndpoint)),
			))

	} else {
		// do not export logs
		logExporter = &NoOpLogExporter{}
	}

	if err != nil {
		return nil, err
	}

	// instantiate log provider with exporter defined above
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName)),
		))

	return loggerProvider, nil
}
