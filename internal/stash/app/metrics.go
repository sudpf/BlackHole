package app

import (
	appconfig "BlackHole/internal/stash/config"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsServer struct {
	server   *http.Server
	listener net.Listener
}

func newMetricsServer(cfg *appconfig.Config) (*metricsServer, error) {
	if cfg == nil || !cfg.MetricsEnabled() {
		return nil, nil
	}

	mux := http.NewServeMux()
	mux.Handle(cfg.MetricsPath(), promhttp.Handler())
	server := &http.Server{
		Addr:    cfg.MetricsListenAddress(),
		Handler: mux,
	}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen metrics: %w", err)
	}

	return &metricsServer{
		server:   server,
		listener: listener,
	}, nil
}

func runMetricsServer(metrics *metricsServer) func(context.Context) error {
	if metrics == nil {
		return nil
	}

	return func(context.Context) error {
		if err := metrics.server.Serve(metrics.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve metrics: %w", err)
		}
		return nil
	}
}

func shutdownMetricsServer(metrics *metricsServer) func(context.Context) error {
	if metrics == nil {
		return nil
	}

	return func(ctx context.Context) error {
		if err := metrics.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown metrics: %w", err)
		}
		return nil
	}
}
