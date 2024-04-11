package logs

import (
	"context"
	"fmt"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogsgrpc"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/stdout/stdoutlogs"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func GetLogProvider(exporterType string, otelExporterEndpoint string, serviceName string) (*sdk.LoggerProvider, error, func()) {

	var logExporter sdk.LogRecordExporter
	var err error

	if exporterType == "stdout" {
		// export logs to stdout of this service
		logExporter, err = stdoutlogs.NewExporter()

	} else if exporterType == "otel" {
		// export logs to an otel collector
		logExporter, err = otlplogs.NewExporter(
			context.Background(),
			otlplogs.WithClient(otlplogsgrpc.NewClient(
				otlplogsgrpc.WithInsecure(),
				otlplogsgrpc.WithEndpoint(otelExporterEndpoint)),
			))

	} else {
		// do not export logs
		logExporter = &NoOpLogExporter{}
	}

	if err != nil {
		return nil, err, nil
	}

	// instantiate log provider with exporter defined above
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName)),
		))

	return loggerProvider, nil, func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error while shutting down Logger provider: %v\n", err)
		}
	}
}
