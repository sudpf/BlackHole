package requestctx

import (
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestResolveTraceIDKeepsValidTraceID(t *testing.T) {
	input := "4bf92f3577b34da6a3ce929d0e0e4736"
	if got := ResolveTraceID(input); got != input {
		t.Fatalf("ResolveTraceID(%q) = %q, want %q", input, got, input)
	}
}

func TestResolveTraceIDGeneratesValidTraceID(t *testing.T) {
	got := ResolveTraceID("trace-123")
	if got == "trace-123" {
		t.Fatalf("ResolveTraceID kept invalid trace ID")
	}
	if _, err := trace.TraceIDFromHex(got); err != nil {
		t.Fatalf("ResolveTraceID generated invalid trace ID %q: %v", got, err)
	}
}

func TestResolveTraceIDRejectsUppercaseTraceID(t *testing.T) {
	input := "4BF92F3577B34DA6A3CE929D0E0E4736"
	got := ResolveTraceID(input)
	if got == input {
		t.Fatalf("ResolveTraceID kept uppercase trace ID")
	}
	if _, err := trace.TraceIDFromHex(got); err != nil {
		t.Fatalf("ResolveTraceID generated invalid trace ID %q: %v", got, err)
	}
}
