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
	Run             func(context.Context) error
	Shutdown        func(context.Context) error
	ShutdownTimeout time.Duration
}

func Run(ctx context.Context, r Runner) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	if r.Run == nil {
		return fmt.Errorf("run function is required")
	}

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	runDone := make(chan error, 1)
	go func() {
		runDone <- r.Run(runCtx)
		close(runDone)
	}()

	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case runErr := <-runDone:
		return runErr
	case <-signalCtx.Done():
		cancelRun()

		if r.Name != "" {
			log.Infof("Shutting down %s", r.Name)
		} else {
			log.Info("Shutting down application")
		}

		shutdownCtx := context.WithoutCancel(ctx)
		var cancel context.CancelFunc
		if r.ShutdownTimeout > 0 {
			shutdownCtx, cancel = context.WithTimeout(shutdownCtx, r.ShutdownTimeout)
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
