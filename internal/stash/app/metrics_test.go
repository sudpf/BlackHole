package app

import (
	appconfig "BlackHole/internal/stash/config"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewMetricsServerDisabled(t *testing.T) {
	server, err := newMetricsServer(&appconfig.Config{})
	if err != nil {
		t.Fatalf("newMetricsServer() error = %v", err)
	}
	if server != nil {
		t.Fatalf("newMetricsServer() = %#v, want nil", server)
	}
}

func TestRunMetricsServerDisabled(t *testing.T) {
	if runner := runMetricsServer(nil); runner != nil {
		t.Fatal("runMetricsServer(nil) returned a runner, want nil")
	}
	if shutdown := shutdownMetricsServer(nil); shutdown != nil {
		t.Fatal("shutdownMetricsServer(nil) returned a shutdown hook, want nil")
	}
}

func TestNewMetricsServerEnabled(t *testing.T) {
	cfg := metricsTestConfig(t)
	server, err := newMetricsServer(cfg)
	if err != nil {
		t.Fatalf("newMetricsServer() error = %v", err)
	}
	defer server.listener.Close()
	if server == nil {
		t.Fatal("newMetricsServer() = nil, want server")
	}
	if server.server.Addr != "127.0.0.1:0" {
		t.Fatalf("server.Addr = %q, want 127.0.0.1:0", server.server.Addr)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	server.server.Handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestRunShutdownsSkipsNil(t *testing.T) {
	called := false
	if err := runShutdowns(context.Background(), nil, func(context.Context) error {
		called = true
		return nil
	}); err != nil {
		t.Fatalf("runShutdowns returned unexpected error: %v", err)
	}
	if !called {
		t.Fatal("non-nil shutdown was not called")
	}
}

func metricsTestConfig(t *testing.T) *appconfig.Config {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "stash.yaml")
	if err := os.WriteFile(path, []byte(`app:
  listen: http://127.0.0.1:8002
  request_timeout: 30s
  shutdown_timeout: 10s
log:
  level: info
  size: 256m
  dir: logs
metrics:
  enabled: true
  listen: http://127.0.0.1:0
  path: /metrics
clusters: []
grace_period: 10s
`), 0o600); err != nil {
		t.Fatalf("write test config: %v", err)
	}

	cfg, err := appconfig.Load(path)
	if err != nil {
		t.Fatalf("load test config: %v", err)
	}
	return cfg
}
