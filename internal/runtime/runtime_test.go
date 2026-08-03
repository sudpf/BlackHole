package runtime

import (
	"context"
	"testing"
	"time"
)

func TestRunRejectsNilContext(t *testing.T) {
	var ctx context.Context
	err := Run(ctx, Runner{
		Run: func(context.Context) error {
			return nil
		},
	})
	if err == nil {
		t.Fatal("Run expected nil context error")
	}
}

func TestRunUsesRootContextAndAllowsGracefulShutdown(t *testing.T) {
	rootCtx, cancelRoot := context.WithCancel(context.Background())
	runStarted := make(chan struct{})
	releaseRun := make(chan struct{})
	shutdownCalled := make(chan error, 1)

	done := make(chan error, 1)
	go func() {
		done <- Run(rootCtx, Runner{
			Name:            "test",
			ShutdownTimeout: time.Second,
			Run: func(ctx context.Context) error {
				close(runStarted)
				<-ctx.Done()
				<-releaseRun
				return nil
			},
			Shutdown: func(ctx context.Context) error {
				shutdownCalled <- ctx.Err()
				close(releaseRun)
				return nil
			},
		})
	}()

	<-runStarted
	cancelRoot()

	select {
	case err := <-shutdownCalled:
		if err != nil {
			t.Fatalf("shutdown context error = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown was not called")
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run error = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not return")
	}
}
