package app

import (
	"BlackHole/api/voidengine/openapi"
	v1handler "BlackHole/api/voidengine/openapi/v1/handler"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/runtime"
	appconfig "BlackHole/internal/voidengine/config"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/service"
	"context"
	"errors"
	"fmt"
	"net/http"
)

func Run(ctx context.Context, cfg *appconfig.Config) (err error) {
	models, err := model.New(ctx, cfg.Database, cfg.LogDir(), cfg.LogSize())
	if err != nil {
		return fmt.Errorf("initialize models: %w", err)
	}
	defer func() {
		if closeErr := models.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close models: %w", closeErr))
		}
	}()

	services := service.New(models)
	handlers := v1handler.New(services)

	listenAddress := cfg.ListenHTTP()
	apiServer, err := openapi.NewHTTPServer(listenAddress, cfg.ApiLogFile(), cfg.LogSize(), cfg.RequestTimeout())
	if err != nil {
		return fmt.Errorf("initialize OpenAPI server: %w", err)
	}
	v1router.RegisterRoutes(apiServer, handlers)
	server := apiServer.HTTPServer()
	if err := runtime.Run(runtime.Runner{
		Name:            "VoidEngine",
		ShutdownTimeout: cfg.ShutdownTimeout(),
		Run: func(context.Context) error {
			if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
				return fmt.Errorf("run OpenAPI: listen and serve: %w", serveErr)
			}
			return nil
		},
		Shutdown: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("shutdown OpenAPI: %w", err)
			}
			return nil
		},
	}); err != nil {
		return err
	}

	return nil
}
