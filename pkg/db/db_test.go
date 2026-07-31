package db

import (
	"BlackHole/pkg/requestctx"
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestLogrusAdapterIncludesTraceID(t *testing.T) {
	var output bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&output)
	logger.SetFormatter(&CustomFormatter{})
	adapter := NewLogrusAdapter(logger)

	ctx := requestctx.WithScope(context.Background(), requestctx.Scope{TraceID: "trace-123"})
	adapter.Trace(ctx, time.Now(), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	var fields map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &fields); err != nil {
		t.Fatalf("unmarshal log: %v: %s", err, output.String())
	}
	if fields["trace_id"] != "trace-123" {
		t.Fatalf("trace_id = %v, want trace-123", fields["trace_id"])
	}
	if _, ok := fields["db"]; ok {
		t.Fatalf("db should not be present: %v", fields)
	}
	if _, ok := fields["elapsed_ms"]; !ok {
		t.Fatalf("elapsed_ms missing: %v", fields)
	}
}

func TestLogrusAdapterMarksStartupWhenTraceIDMissing(t *testing.T) {
	var output bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&output)
	logger.SetFormatter(&CustomFormatter{})
	adapter := NewLogrusAdapter(logger)

	adapter.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "SELECT 1", 1
	}, nil)

	var fields map[string]interface{}
	if err := json.Unmarshal(output.Bytes(), &fields); err != nil {
		t.Fatalf("unmarshal log: %v: %s", err, output.String())
	}
	if fields["trace_id"] != "system" {
		t.Fatalf("trace_id = %v, want system", fields["trace_id"])
	}
}
