# otelhandlers

While developing a system that makes use of logs, metrics, and traces reporting with OpenTelemetry, it is generally 
useful to switch the reporting destination for each telemetry type between `stdout`, a `collector` instance, and 
silencing them altogether. This library allows for the exporter for each of these telemetry types to be determined
by a single flag (which in the host application can be set via an environment variable or configuration parameter).

## usage

To configure telemetry for some type, OpenTelemetry requires that you first define an `exporter`, which is then
passed to a `provider`. The configured `provider` instance is then set globally to intercept and format messages
before being sent to their reporting destination. 

Each module of this library provides two functions: `GetExporter` and `GetProvider`. The `GetExporter` function
is only needed if you are configuring your own `provider` instance. Otherwise, the `GetProvider` function is 
sufficient on its own. 

Below, usage instructions are provided for each of the 3 telemetry types:

### logs

I have only configured OTel logs alongside the `zap` logging framework, but the procedure for other frameworks
or the default `log` library should be analogous to the following:
```go
import (
    "github.com/agoda-com/opentelemetry-go/otelzap"
    "github.com/bengetch/otelhandlers/logs"
    "go.uber.org/zap"
)

lp, lpErr := logshandler.GetLogProvider(os.Getenv("LOGS_EXPORTER"), "my-service")
logger := zap.New(otelzap.NewOtelCore(lp))
zap.ReplaceGlobals(logger)
```

In the above code, the environment variable `LOGS_EXPORTER` can be set to either `stdout`, `otel,` or `noop`, 
and a corresponding provider instance will be returned. Logs can then be made via:

```go
otelzap.Ctx(<some-context>).Info("hello from my-service")
```


### metrics

To configure OTel metrics:
```go
import (
    "github.com/bengetch/otelhandlers/metrics"
)

mp, mpErr := metricshandler.GetMeterProvider(os.Getenv("METRICS_EXPORTER"), "my-service")
otel.SetMeterProvider(mp)
```

Metrics can then be recorded in the normal ways, and output will be sent to whichever reporting destination you
configured via the `METRICS_EXPORTER` environment variable above.


### traces

To configure OTel traces:
```go
tp, tpErr := traceshandler.GetTracerProvider(os.Getenv("TRACES_EXPORTER"), ServiceName)
otel.SetTracerProvider(tp)
textPropagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
otel.SetTextMapPropagator(textPropagator)
```

Traces are a unique case because we typically want to propagate trace context across API endpoints, in order to enable
visualization tooling like flame graphs. To accomplish this, you just need to pass the current context along with any
HTTP requests, and extract the request context from any API endpoint handlers. 