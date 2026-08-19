package app

import (
	"BlackHole/internal/runtime"
	appconfig "BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service"
	"context"
	"errors"
	"fmt"
)

func Run(ctx context.Context, cfg *appconfig.Config) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	stash, err := service.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("initialize service: %w", err)
	}
	metricsServer, err := newMetricsServer(cfg)
	if err != nil {
		return fmt.Errorf("initialize metrics server: %w", err)
	}

	if err := runtime.Run(ctx, runtime.Runner{
		Name:            "Stash",
		ShutdownTimeout: cfg.ShutdownTimeout(),
		Run: func(ctx context.Context) error {
			runBackground(ctx, runMetricsServer(metricsServer))
			return runStashService(stash)(ctx)
		},
		Shutdown: func(ctx context.Context) error {
			return runShutdowns(ctx,
				shutdownMetricsServer(metricsServer),
				stopStashService(stash),
			)
		},
	}); err != nil {
		return err
	}

	return nil
}

func runBackground(ctx context.Context, runner func(context.Context) error) {
	if runner == nil {
		return
	}
	go func() {
		_ = runner(ctx)
	}()
}

func runShutdowns(ctx context.Context, shutdowns ...func(context.Context) error) error {
	var err error
	for _, shutdown := range shutdowns {
		if shutdown != nil {
			err = errors.Join(err, shutdown(ctx))
		}
	}
	return err
}

func runStashService(stash *service.Stash) func(context.Context) error {
	return func(ctx context.Context) error {
		if err := stash.Run(ctx); err != nil {
			return fmt.Errorf("run service: %w", err)
		}
		return nil
	}
}

func stopStashService(stash *service.Stash) func(context.Context) error {
	return func(ctx context.Context) error {
		if err := stash.Stop(ctx); err != nil {
			return fmt.Errorf("stop service: %w", err)
		}
		return nil
	}
}
