package app

import (
	"BlackHole/internal/runtime"
	appconfig "BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service"
	"context"
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

	if err := runtime.Run(ctx, runtime.Runner{
		Name:            "Stash",
		ShutdownTimeout: cfg.ShutdownTimeout(),
		Run: func(ctx context.Context) error {
			if err := stash.Run(ctx); err != nil {
				return fmt.Errorf("run service: %w", err)
			}
			return nil
		},
		Shutdown: func(ctx context.Context) error {
			if err := stash.Stop(ctx); err != nil {
				return fmt.Errorf("stop service: %w", err)
			}
			return nil
		},
	}); err != nil {
		return err
	}

	return nil
}
