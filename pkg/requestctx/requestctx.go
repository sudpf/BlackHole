package requestctx

import (
	"context"
	"strings"

	"github.com/google/uuid"
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
	if validTraceID(value) {
		return value
	}

	return uuid.NewString()
}

func validTraceID(value string) bool {
	if len(value) == 0 || len(value) > 64 {
		return false
	}

	for _, char := range value {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == '.' {
			continue
		}
		return false
	}

	return true
}
