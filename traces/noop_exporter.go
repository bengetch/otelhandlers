package traces

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type NoOpSpanExporter struct{}

func (e *NoOpSpanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	return nil
}

func (e *NoOpSpanExporter) Shutdown(ctx context.Context) error { return nil }
