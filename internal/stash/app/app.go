package app

import (
	"BlackHole/internal/stash/service"
	"BlackHole/pkg/config"
	"fmt"
)

func Run(cfg *config.StashConfig) error {
	if err := service.Init(cfg); err != nil {
		return fmt.Errorf("initialize service: %w", err)
	}

	service.Run()
	return nil
}
