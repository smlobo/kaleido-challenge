package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// A refusal is visible in the trace: its code as error.type, and its title
// and detail in a "problem" event. Only a 5xx marks the span as an error.
func TestProblemIsRecordedOnTheSpan(t *testing.T) {
	rec := tracetest.NewSpanRecorder()
	tracer := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec)).Tracer("test")

	for _, tc := range []struct {
		status    int
		code      string
		wantError bool
	}{
		{http.StatusUnprocessableEntity, "pass_expired", false},
		{http.StatusInternalServerError, "internal", true},
	} {
		ctx, span := tracer.Start(t.Context(), "POST /tickets/{id}/check-in")
		r := httptest.NewRequest(http.MethodPost, "/tickets/x/check-in", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		writeProblem(w, r, tc.status, tc.code, "Title", "Detail")
		span.End()

		if w.Code != tc.status {
			t.Fatalf("%s: answered %d", tc.code, w.Code)
		}
		got := rec.Ended()[len(rec.Ended())-1]
		attrs := map[string]string{}
		for _, a := range got.Attributes() {
			attrs[string(a.Key)] = a.Value.Emit()
		}
		if attrs["error.type"] != tc.code {
			t.Errorf("%s: error.type = %q", tc.code, attrs["error.type"])
		}
		events := got.Events()
		if len(events) != 1 || events[0].Name != "problem" {
			t.Fatalf("%s: events %+v", tc.code, events)
		}
		ev := map[string]string{}
		for _, a := range events[0].Attributes {
			ev[string(a.Key)] = a.Value.Emit()
		}
		if ev["problem.code"] != tc.code || ev["problem.title"] != "Title" || ev["problem.detail"] != "Detail" {
			t.Errorf("%s: problem event %v", tc.code, ev)
		}
		if isError := got.Status().Code == codes.Error; isError != tc.wantError {
			t.Errorf("%s: span error = %v, want %v", tc.code, isError, tc.wantError)
		}
	}
}
