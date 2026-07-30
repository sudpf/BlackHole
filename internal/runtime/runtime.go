package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Runner struct {
	Name            string
	Run             func() error
	Shutdown        func(context.Context) error
	ShutdownTimeout time.Duration
}

func Run(r Runner) error {
	if r.Run == nil {
		return fmt.Errorf("run function is required")
	}

	runDone := make(chan error, 1)
	go func() {
		runDone <- r.Run()
		close(runDone)
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case runErr := <-runDone:
		return runErr
	case <-signalCtx.Done():
		if r.Name != "" {
			log.Infof("Shutting down %s", r.Name)
		} else {
			log.Info("Shutting down application")
		}

		shutdownCtx := context.Background()
		var cancel context.CancelFunc
		if r.ShutdownTimeout > 0 {
			shutdownCtx, cancel = context.WithTimeout(context.Background(), r.ShutdownTimeout)
			defer cancel()
		}

		if r.Shutdown != nil {
			if err := r.Shutdown(shutdownCtx); err != nil {
				return err
			}
		}

		if r.ShutdownTimeout > 0 {
			select {
			case runErr := <-runDone:
				return runErr
			case <-shutdownCtx.Done():
				return fmt.Errorf("shutdown %s: %w", r.Name, shutdownCtx.Err())
			}
		}

		return <-runDone
	}
}
