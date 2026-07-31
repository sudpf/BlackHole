package requestctx

import (
	"context"
	"crypto/rand"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

const HeaderTraceID = "X-Trace-ID"

type scopeKey struct{}

type Scope struct {
	TraceID  string
	Language string
	ClientIP string
}

func WithScope(ctx context.Context, scope Scope) context.Context {
	return context.WithValue(ctx, scopeKey{}, scope)
}

func TraceID(ctx context.Context) string {
	scope, ok := ctx.Value(scopeKey{}).(Scope)
	if !ok {
		return ""
	}

	return scope.TraceID
}

func Language(ctx context.Context) string {
	scope, ok := ctx.Value(scopeKey{}).(Scope)
	if !ok {
		return ""
	}

	return scope.Language
}

func ClientIP(ctx context.Context) string {
	scope, ok := ctx.Value(scopeKey{}).(Scope)
	if !ok {
		return ""
	}

	return scope.ClientIP
}

func ResolveTraceID(value string) string {
	value = strings.TrimSpace(value)
	if _, err := trace.TraceIDFromHex(value); err == nil {
		return value
	}

	traceID := trace.TraceID{}
	if _, err := rand.Read(traceID[:]); err != nil {
		panic("generate trace id: " + err.Error())
	}
	return traceID.String()
}
