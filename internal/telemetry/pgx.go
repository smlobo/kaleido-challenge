package telemetry

import (
	"context"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// PGXTracer gives each Postgres query a span, as a child of the span in the
// query's context. Queries outside a trace (the outbox loops polling for
// work, health checks) get none.
//
// Set it as pgx.ConnConfig.Tracer.
type PGXTracer struct{}

var _ pgx.QueryTracer = PGXTracer{}

type pgxSpanKey struct{}

// TraceQueryStart starts the query's span.
func (PGXTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if !InTrace(ctx) {
		return ctx
	}
	sql := compact(data.SQL)
	ctx, span := Tracer().Start(ctx, sqlSpanName(sql), trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system.name", "postgresql"),
			attribute.String("db.query.text", sql),
		))
	return context.WithValue(ctx, pgxSpanKey{}, span)
}

// TraceQueryEnd ends it, recording the rows affected or the error.
func (PGXTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span, ok := ctx.Value(pgxSpanKey{}).(trace.Span)
	if !ok {
		return
	}
	if data.Err != nil {
		span.RecordError(data.Err)
		span.SetStatus(codes.Error, data.Err.Error())
	} else {
		span.SetAttributes(attribute.Int64("db.response.rows_affected", data.CommandTag.RowsAffected()))
	}
	span.End()
}

var (
	spaces     = regexp.MustCompile(`\s+`)
	firstTable = regexp.MustCompile(`(?i)\b(?:from|into|update|join)\s+([a-z_][a-z0-9_.]*)`)
	innermost  = regexp.MustCompile(`\([^()]*\)`)
)

// compact puts a query on one line, for reading in a trace viewer.
func compact(sql string) string { return strings.TrimSpace(spaces.ReplaceAllString(sql, " ")) }

// sqlSpanName names a span by the statement and its first table, e.g.
// "UPDATE transactions" or "SELECT tickets"; "BEGIN" and "COMMIT" stay as
// they are. Parenthesised parts (subqueries, column lists) are skipped, so a
// subquery in the select list does not name the query.
func sqlSpanName(sql string) string {
	verb, _, _ := strings.Cut(sql, " ")
	verb = strings.ToUpper(verb)
	outer := sql
	for {
		stripped := innermost.ReplaceAllString(outer, " ")
		if stripped == outer {
			break
		}
		outer = stripped
	}
	if m := firstTable.FindStringSubmatch(outer); m != nil {
		return verb + " " + strings.ToLower(m[1])
	}
	return verb
}
