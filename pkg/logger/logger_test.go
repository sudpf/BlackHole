package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"

	"BlackHole/pkg/requestctx"
)

func TestRotatingWriterRejectsLessThanOneMB(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "test.log")
	if _, err := RotatingWriter(filename, "10240"); err == nil {
		t.Fatal(`RotatingWriter("test.log", "10240") expected error`)
	}
}

func TestRotatingWriterReturnsPrepareLogFileError(t *testing.T) {
	dir := t.TempDir()
	notDir := filepath.Join(dir, "not-dir")
	if err := os.WriteFile(notDir, []byte("file"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	if _, err := RotatingWriter(filepath.Join(notDir, "test.log"), "1m"); err == nil {
		t.Fatal("RotatingWriter expected prepare log file error")
	}
}

func TestJSONFormatterOmitsCaller(t *testing.T) {
	var output bytes.Buffer
	logger := log.New()
	logger.SetFormatter(JSONFormatter())
	logger.SetOutput(&output)
	logger.SetReportCaller(true)

	logger.Info("test")

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log entry: %v", err)
	}
	if _, ok := entry["caller"]; ok {
		t.Fatalf("unexpected caller field: %v", entry["caller"])
	}
	if _, ok := entry["file"]; !ok {
		t.Fatal("missing file field")
	}
}

func TestFromContextOnlyIncludesTraceID(t *testing.T) {
	ctx := requestctx.WithScope(context.Background(), requestctx.Scope{
		TraceID:  "trace-1",
		ClientIP: "127.0.0.1",
		Language: "zh-CN",
	})

	entry := FromContext(ctx)
	if entry.Data["trace_id"] != "trace-1" {
		t.Fatalf("trace_id = %v, want trace-1", entry.Data["trace_id"])
	}
	if _, ok := entry.Data["client_ip"]; ok {
		t.Fatalf("unexpected client_ip field: %v", entry.Data["client_ip"])
	}
	if _, ok := entry.Data["language"]; ok {
		t.Fatalf("unexpected language field: %v", entry.Data["language"])
	}
}
