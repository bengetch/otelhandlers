package sdk

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func CleanupLoggerProvider(lp *log.LoggerProvider) {

	ctx := context.Background()
	if lpShutdownErr := lp.Shutdown(ctx); lpShutdownErr != nil {
		slog.ErrorContext(
			context.Background(),
			fmt.Sprintf("error while shutting down logger provider: %v\n", lpShutdownErr),
		)
	}
}

func CleanupTracerProvider(tp *trace.TracerProvider) {

	ctx := context.Background()
	if tpShutdownErr := tp.Shutdown(ctx); tpShutdownErr != nil {
		slog.ErrorContext(
			context.Background(),
			fmt.Sprintf("error while shutting down tracer provider: %v\n", tpShutdownErr),
		)
	}
}

func CleanupMeterProvider(mp *metric.MeterProvider) {

	ctx := context.Background()
	if mpShutdownErr := mp.Shutdown(ctx); mpShutdownErr != nil {
		slog.ErrorContext(
			context.Background(),
			fmt.Sprintf("error while shutting down meter provider: %v\n", mpShutdownErr),
		)
	}
}

func CleanupTelemetryProviders(lp *log.LoggerProvider, tp *trace.TracerProvider, mp *metric.MeterProvider) {
	/*
		call shutdown function on all telemetry provider types
	*/

	if lp != nil {
		CleanupLoggerProvider(lp)
	}

	if tp != nil {
		CleanupTracerProvider(tp)
	}

	if mp != nil {
		CleanupMeterProvider(mp)
	}
}
