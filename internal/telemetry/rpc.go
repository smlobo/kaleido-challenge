package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Transport wraps an http.RoundTripper for JSON-RPC calls to the node. Inside
// a trace, each call gets a client span named after its JSON-RPC method (e.g.
// "eth_sendRawTransaction") and carries the trace headers, so the node's own
// spans become its children. Outside a trace it adds nothing.
//
// Use it as the go-ethereum rpc client's HTTP transport.
type Transport struct {
	Base http.RoundTripper // http.DefaultTransport if nil
}

// RoundTrip implements http.RoundTripper.
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	ctx := req.Context()
	if !InTrace(ctx) {
		return base.RoundTrip(req)
	}

	// A RoundTripper must not change the request it is given: work on a copy,
	// whose body (read here for the method name) and headers are its own.
	req = req.Clone(ctx)
	method, err := rpcMethod(req)
	if err != nil {
		return nil, err
	}
	ctx, span := Tracer().Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("rpc.system", "jsonrpc"),
			attribute.String("rpc.method", method),
			attribute.String("server.address", req.URL.Host),
		))
	defer span.End()

	req = req.WithContext(ctx)
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := base.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, resp.Status)
	}
	return resp, nil
}

// rpcMethod reads the JSON-RPC method from the request body and puts the body
// back. A batch is named by its methods, joined.
func rpcMethod(req *http.Request) (string, error) {
	if req.Body == nil {
		return "jsonrpc", nil
	}
	body, err := io.ReadAll(req.Body)
	_ = req.Body.Close()
	if err != nil {
		return "", fmt.Errorf("read JSON-RPC request: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	type call struct {
		Method string `json:"method"`
	}
	var one call
	if json.Unmarshal(body, &one) == nil && one.Method != "" {
		return one.Method, nil
	}
	var batch []call
	if json.Unmarshal(body, &batch) == nil && len(batch) > 0 {
		names := make([]string, len(batch))
		for i, c := range batch {
			names[i] = c.Method
		}
		return "batch " + strings.Join(names, ","), nil
	}
	return "jsonrpc", nil
}
