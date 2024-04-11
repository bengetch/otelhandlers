package logs

import (
	"context"
	sdk "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
)

type NoOpLogExporter struct{}

func (e *NoOpLogExporter) Shutdown(ctx context.Context) error { return nil }

func (e *NoOpLogExporter) Export(ctx context.Context, batch []sdk.ReadableLogRecord) error {
	return nil
}
