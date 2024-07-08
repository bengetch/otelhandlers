package logs

import (
	"context"

	"go.opentelemetry.io/otel/sdk/log"
)

type NoOpLogExporter struct{}

func (e NoOpLogExporter) Shutdown(ctx context.Context) error { return nil }

func (e NoOpLogExporter) Export(ctx context.Context, records []log.Record) error {
	return nil
}

func (e NoOpLogExporter) ForceFlush(ctx context.Context) error { return nil }
