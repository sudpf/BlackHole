package app

import (
	"BlackHole/api/voidengine/openapi"
	v1handler "BlackHole/api/voidengine/openapi/v1/handler"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/config"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func Run(cfg *config.VoidEngineConfig) (err error) {
	models, err := model.New(cfg.Database)
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
	v1router.RegisterRoutes(handlers)

	listenAddress := cfg.ListenHTTP()
	openapi.InitApi(listenAddress)

	server := openapi.NewServer(listenAddress)
	serverErr := make(chan error, 1)
	go func() {
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("listen and serve: %w", serveErr)
		}
		close(serverErr)
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case serveErr := <-serverErr:
		if serveErr != nil {
			return fmt.Errorf("run OpenAPI: %w", serveErr)
		}
	case <-signalCtx.Done():
		log.Info("Shutting down VoidEngine")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout())
		defer cancel()

		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			return fmt.Errorf("shutdown OpenAPI: %w", shutdownErr)
		}

		if serveErr := <-serverErr; serveErr != nil {
			return fmt.Errorf("run OpenAPI: %w", serveErr)
		}
	}

	return nil
}
