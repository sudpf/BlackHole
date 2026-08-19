package service

import (
	"BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service/handler"
	"context"
	"errors"
	"fmt"

	zeroservice "github.com/zeromicro/go-zero/core/service"
)

type Stash struct {
	group    *zeroservice.ServiceGroup
	handlers []*handler.MessageHandler
}

func New(ctx context.Context, cfg *config.Config) (_ *Stash, err error) {
	return NewBuilder().Build(ctx, cfg)
}

func (s *Stash) Run(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil || s.group == nil {
		return fmt.Errorf("service group not initialized")
	}

	s.group.Start()
	return nil
}

func (s *Stash) Stop(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}
	if s == nil || s.group == nil {
		return nil
	}

	s.group.Stop()

	return closeHandlers(ctx, s.handlers)
}

func closeHandlers(ctx context.Context, handlers []*handler.MessageHandler) error {
	var err error
	for _, handle := range handlers {
		err = errors.Join(err, handle.Close(ctx))
	}

	return err
}
