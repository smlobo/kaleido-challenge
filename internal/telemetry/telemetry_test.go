package telemetry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// record installs a tracer provider that keeps finished spans in memory.
func record(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	rec := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })
	return rec
}

func TestSQLSpanName(t *testing.T) {
	for sql, want := range map[string]string{
		"SELECT id FROM transactions WHERE id = $1":            "SELECT transactions",
		"UPDATE tickets t SET x = 1 FROM resale_actions a":     "UPDATE tickets",
		"INSERT INTO transactions (a) VALUES ($1)":             "INSERT transactions",
		"select pg_advisory_lock($1)":                          "SELECT",
		"begin":                                                "BEGIN",
		"SAVEPOINT sp":                                         "SAVEPOINT",
		"WITH x AS (SELECT 1 FROM events) SELECT * FROM x":     "WITH x",
		"DELETE FROM api_tokens WHERE token_hash = $1":         "DELETE api_tokens",
		"SELECT t.id FROM tickets t JOIN events e ON e.id = 1": "SELECT tickets",
		"SELECT t.id, EXISTS (SELECT 1 FROM resale_actions a WHERE (a.x = 1)) FROM tickets t FOR UPDATE OF t": "SELECT tickets",
	} {
		if got := sqlSpanName(compact(sql)); got != want {
			t.Errorf("%q: got %q, want %q", sql, got, want)
		}
	}
}

func TestInjectExtractRoundTrip(t *testing.T) {
	rec := record(t)
	if got := Inject(context.Background()); got != "" {
		t.Fatalf("outside a trace Inject = %q, want empty", got)
	}
	ctx, span := Tracer().Start(context.Background(), "request")
	tp := Inject(ctx)
	span.End()
	if !strings.HasPrefix(tp, "00-"+span.SpanContext().TraceID().String()+"-") {
		t.Fatalf("traceparent %q does not carry trace %s", tp, span.SpanContext().TraceID())
	}

	// Later, in a loop: the saved traceparent makes the loop's span a child.
	_, child := Child(Extract(context.Background(), tp), "outbox receipt")
	child.End()
	spans := rec.Ended()
	if len(spans) != 2 || spans[1].Parent().SpanID() != span.SpanContext().SpanID() ||
		spans[1].SpanContext().TraceID() != span.SpanContext().TraceID() {
		t.Fatalf("loop span is not a child of the request: %+v", spans)
	}

	// No saved context: no span at all.
	_, none := Child(Extract(context.Background(), ""), "outbox receipt")
	none.End()
	if len(rec.Ended()) != 2 {
		t.Fatal("a span was started outside any trace")
	}
}

func TestTransportNamesSpanAndPropagates(t *testing.T) {
	otel.SetTextMapPropagator(propagator)
	rec := record(t)
	var got http.Header
	var body string
	node := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":"0x1"}`)
	}))
	defer node.Close()
	client := &http.Client{Transport: Transport{}}
	call := `{"jsonrpc":"2.0","id":1,"method":"eth_sendRawTransaction","params":["0x01"]}`

	// Outside a trace: no headers, no span.
	if _, err := client.Post(node.URL, "application/json", strings.NewReader(call)); err != nil {
		t.Fatal(err)
	}
	if got.Get("traceparent") != "" || got.Get("b3") != "" || len(rec.Ended()) != 0 {
		t.Fatalf("untraced call was traced: %v %d", got, len(rec.Ended()))
	}

	ctx, parent := Tracer().Start(context.Background(), "request")
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, node.URL, strings.NewReader(call))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	parent.End()

	if body != call {
		t.Fatalf("the node got %q, want the original body", body)
	}
	spans := rec.Ended()
	if len(spans) != 2 || spans[0].Name() != "eth_sendRawTransaction" || spans[0].SpanKind() != trace.SpanKindClient {
		t.Fatalf("spans: %+v", spans)
	}
	rpcSpan := spans[0].SpanContext()
	if !strings.Contains(got.Get("traceparent"), rpcSpan.SpanID().String()) {
		t.Fatalf("traceparent %q does not name the RPC span %s", got.Get("traceparent"), rpcSpan.SpanID())
	}
	// Besu reads b3, not traceparent: trace-span-sampled.
	if want := rpcSpan.TraceID().String() + "-" + rpcSpan.SpanID().String() + "-1"; got.Get("b3") != want {
		t.Fatalf("b3 = %q, want %q", got.Get("b3"), want)
	}
}

func TestRPCMethodBatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`[{"method":"eth_chainId"},{"method":"eth_blockNumber"}]`))
	if m, err := rpcMethod(req); err != nil || m != "batch eth_chainId,eth_blockNumber" {
		t.Fatalf("got %q %v", m, err)
	}
}

func TestPGXTracerOnlyInsideATrace(t *testing.T) {
	rec := record(t)
	var tr PGXTracer
	data := pgx.TraceQueryStartData{SQL: "UPDATE transactions\n\t SET status = $2 WHERE id = $1"}
	tag := pgconn.NewCommandTag("UPDATE 1")

	ctx := tr.TraceQueryStart(context.Background(), nil, data)
	tr.TraceQueryEnd(ctx, nil, pgx.TraceQueryEndData{CommandTag: tag})
	if len(rec.Ended()) != 0 {
		t.Fatal("a query outside any trace got a span")
	}

	ctx, parent := Tracer().Start(context.Background(), "request")
	qctx := tr.TraceQueryStart(ctx, nil, data)
	tr.TraceQueryEnd(qctx, nil, pgx.TraceQueryEndData{CommandTag: tag})
	parent.End()
	spans := rec.Ended()
	if len(spans) != 2 || spans[0].Name() != "UPDATE transactions" {
		t.Fatalf("spans: %+v", spans)
	}
	attrs := map[string]string{}
	for _, a := range spans[0].Attributes() {
		attrs[string(a.Key)] = a.Value.Emit()
	}
	if attrs["db.query.text"] != "UPDATE transactions SET status = $2 WHERE id = $1" || attrs["db.response.rows_affected"] != "1" {
		t.Fatalf("attributes: %v", attrs)
	}
}
