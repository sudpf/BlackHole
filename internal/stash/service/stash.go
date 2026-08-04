package service

import (
	"BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service/filter"
	"BlackHole/internal/stash/service/handler"
	"BlackHole/internal/stash/service/input"
	"BlackHole/internal/stash/service/output"
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
)

type Stash struct {
	group    *service.ServiceGroup
	handlers []*handler.MessageHandler
}

func New(ctx context.Context, cfg *config.Config) (_ *Stash, err error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	proc.SetTimeToForceQuit(cfg.GracePeriod)
	stash := &Stash{group: service.NewServiceGroup()}
	var initialized []service.Service
	defer func() {
		if err == nil {
			return
		}
		for _, svc := range initialized {
			svc.Stop()
		}
		for _, handle := range stash.handlers {
			_ = handle.Close(ctx)
		}
	}()

	for i, cluster := range cfg.Clusters {
		if cluster == nil {
			return nil, fmt.Errorf("cluster %d is nil", i)
		}
		if cluster.Input == nil {
			return nil, fmt.Errorf("cluster %d input is nil", i)
		}

		filters := filter.CreateFilters(cluster)

		writers, err := output.NewWriters(ctx, cluster.Output)
		if err != nil {
			return nil, fmt.Errorf("initialize cluster %d writers: %w", i, err)
		}

		handle := handler.NewHandler()
		handle.AddFilters(filters...)
		handle.AddWriters(writers...)
		stash.handlers = append(stash.handlers, handle)

		if cluster.Input.Kafka != nil {
			for j, k := range input.ToKqConf(cluster.Input.Kafka) {
				queue, err := kq.NewQueue(k, handle)
				if err != nil {
					return nil, fmt.Errorf("initialize cluster %d kafka queue %d: %w", i, j, err)
				}
				initialized = append(initialized, queue)
				stash.group.Add(queue)
			}
		}

		if cluster.Input.Syslogs != nil {
			for j, s := range cluster.Input.Syslogs {
				if s == nil {
					return nil, fmt.Errorf("cluster %d syslog input %d is nil", i, j)
				}
				syslogService, err := input.NewSyslogService(ctx, s, handle)
				if err != nil {
					return nil, fmt.Errorf("initialize cluster %d syslog input %d: %w", i, j, err)
				}
				initialized = append(initialized, syslogService)
				stash.group.Add(syslogService)
			}
		}
	}

	return stash, nil
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

	var err error
	for _, handle := range s.handlers {
		err = errors.Join(err, handle.Close(ctx))
	}

	return err
}
