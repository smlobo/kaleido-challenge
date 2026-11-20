// Package telemetry sets up distributed tracing with OpenTelemetry, and holds
// the small pieces that put the service's work into traces: Postgres queries
// (PGXTracer), JSON-RPC calls to the node (Transport), and the trace context
// saved with each outbox transaction (Inject, Extract).
//
// Tracing is off unless OTEL_EXPORTER_OTLP_ENDPOINT (or
// OTEL_EXPORTER_OTLP_TRACES_ENDPOINT) is set. Off, every span is a no-op, so
// go run and the tests behave exactly as before. The exporter reads the
// standard OTEL_* variables itself (endpoint, headers, timeout), and
// OTEL_SERVICE_NAME names the service ("ticketing" by default).
package telemetry

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ScopeName names the tracer for the service's own spans.
const ScopeName = "github.com/smlobo/kaleido-challenge"

// Tracer returns the service's tracer. It follows whatever provider Setup
// installed, and is a no-op before Setup or without an endpoint.
func Tracer() trace.Tracer { return otel.Tracer(ScopeName) }

// propagator writes and reads trace context on the wire. W3C trace context is
// the standard; B3 (single "b3" header) is added because Besu reads B3 and
// Jaeger headers but not "traceparent", so it is how Besu's JSON-RPC spans
// join our traces.
var propagator = propagation.NewCompositeTextMapPropagator(
	propagation.TraceContext{},
	propagation.Baggage{},
	b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)),
)

// Enabled reports whether an OTLP endpoint is configured.
func Enabled() bool {
	return os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" || os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""
}

// Setup installs the propagator and, if an endpoint is configured, a tracer
// provider that batches spans to it over OTLP/HTTP. The returned function
// flushes and stops the provider; call it on shutdown.
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	otel.SetTextMapPropagator(propagator)
	if !Enabled() {
		return func(context.Context) error { return nil }, nil
	}
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}
	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

// InstanceID names this run of the service in traces (service.instance.id):
// a random UUID per process start, so several instances behind one load
// balancer, or one instance before and after a restart, can be told apart.
var InstanceID = uuid.NewString()

// newResource describes the process that makes the spans: the service's name
// ("ticketing"), InstanceID, the host (in a container, its hostname is the
// container ID) and the process. OTEL_SERVICE_NAME and OTEL_RESOURCE_ATTRIBUTES
// override any of these: the environment detector runs last, and later
// detectors win.
func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.String("service.name", "ticketing"),
			attribute.String("service.instance.id", InstanceID),
		),
		resource.WithHost(),
		resource.WithProcessPID(),
		resource.WithFromEnv(),
	)
	// A detector that only partly worked still gives a usable resource.
	if err != nil && res == nil {
		return nil, err
	}
	return res, nil
}

// Inject returns the W3C traceparent of the span in ctx, or "" if ctx has no
// sampled span. The outbox saves it with each transaction.
func Inject(ctx context.Context) string {
	if !trace.SpanContextFromContext(ctx).IsSampled() {
		return ""
	}
	c := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, c)
	return c.Get("traceparent")
}

// Extract returns ctx with the remote span a saved traceparent names as its
// parent. An empty or malformed traceparent leaves ctx unchanged.
func Extract(ctx context.Context, traceparent string) context.Context {
	if traceparent == "" {
		return ctx
	}
	return propagation.TraceContext{}.Extract(ctx, propagation.MapCarrier{"traceparent": traceparent})
}

// InTrace reports whether ctx carries a span to hang child spans on. The
// pieces here create spans only inside a trace, so background polling with
// nothing to do never starts traces of its own.
func InTrace(ctx context.Context) bool {
	return trace.SpanContextFromContext(ctx).IsValid()
}

// Child starts a span under the one in ctx. Outside a trace it starts
// nothing and returns ctx with a no-op span, so callers need no checks.
func Child(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	if !InTrace(ctx) {
		return ctx, trace.SpanFromContext(ctx)
	}
	return Tracer().Start(ctx, name, trace.WithAttributes(attrs...))
}

// EndSpan ends a span, marking it failed if err is not nil.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
