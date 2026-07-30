package app

import (
	"BlackHole/internal/runtime"
	"BlackHole/internal/stash/service"
	"BlackHole/pkg/config"
	"context"
	"fmt"
)

func Run(cfg *config.StashConfig) error {
	stash, err := service.New(cfg)
	if err != nil {
		return fmt.Errorf("initialize service: %w", err)
	}

	if err := runtime.Run(runtime.Runner{
		Name: "Stash",
		Run: func() error {
			if err := stash.Run(); err != nil {
				return fmt.Errorf("run service: %w", err)
			}
			return nil
		},
		Shutdown: func(context.Context) error {
			if err := stash.Stop(); err != nil {
				return fmt.Errorf("stop service: %w", err)
			}
			return nil
		},
	}); err != nil {
		return err
	}

	return nil
}
