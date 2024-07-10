package sdk

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bengetch/otelhandlers/traces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetupTraces(exporterType string, serviceName string) *trace.TracerProvider {
	/*
		configure tracer provider instance, which is responsible for exporting traces to the backend
		indicated by exporterType. the text map propagator configured here also ensures that trace
		context is propagated correctly across API calls
	*/

	tp, tpErr := traces.GetProvider(exporterType, serviceName)
	if tpErr != nil {
		slog.ErrorContext(context.Background(), fmt.Sprintf("Failed to get tracer provider: %v\n", tpErr))
	} else {
		otel.SetTracerProvider(tp)
		textPropagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
		otel.SetTextMapPropagator(textPropagator)
	}

	return tp
}
