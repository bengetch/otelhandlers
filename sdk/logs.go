package sdk

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bengetch/otelhandlers/logs"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
)

func SetupLogs(exporterType string, serviceName string) *log.LoggerProvider {
	/*
		configure logger provider instance, which is responsible for exporting logs to the backend
		indicated by exoorterType
	*/

	lp, lpErr := logs.GetProvider(exporterType, serviceName)
	if lpErr != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("Failed to get log provider: %v\n", lpErr))
	} else {
		logger := otelslog.NewLogger(
			serviceName,
			otelslog.WithLoggerProvider(lp),
		)
		slog.SetDefault(logger)
	}

	return lp
}
