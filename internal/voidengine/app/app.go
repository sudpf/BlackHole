package app

import (
	"BlackHole/api/voidengine/openapi"
	v1handler "BlackHole/api/voidengine/openapi/v1/handler"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/config"
	"fmt"
)

func Run(cfg *config.VoidEngineConfig) error {
	models, err := model.New(cfg.Database)
	if err != nil {
		return fmt.Errorf("initialize models: %w", err)
	}

	services := service.New(models)
	handlers := v1handler.New(services)
	v1router.RegisterRoutes(handlers)

	listenAddress := cfg.ListenHTTP()
	openapi.InitApi(listenAddress)
	if err := openapi.Run(listenAddress); err != nil {
		return fmt.Errorf("run OpenAPI: %w", err)
	}

	return nil
}
