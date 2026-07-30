package service

import (
	"BlackHole/internal/stash/service/filter"
	"BlackHole/internal/stash/service/handler"
	"BlackHole/internal/stash/service/input"
	"BlackHole/internal/stash/service/output"
	"BlackHole/pkg/config"
	"fmt"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
)

var (
	group *service.ServiceGroup
)

func Init(cfg *config.StashConfig) (err error) {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	proc.SetTimeToForceQuit(cfg.GracePeriod)
	group = service.NewServiceGroup()
	var initialized []service.Service
	defer func() {
		if err == nil {
			return
		}
		for _, svc := range initialized {
			svc.Stop()
		}
	}()

	for i, cluster := range cfg.Clusters {
		if cluster == nil {
			return fmt.Errorf("cluster %d is nil", i)
		}
		if cluster.Input == nil {
			return fmt.Errorf("cluster %d input is nil", i)
		}

		filters := filter.CreateFilters(cluster)

		writers, err := output.NewWriters(cluster.Output)
		if err != nil {
			return fmt.Errorf("initialize cluster %d writers: %w", i, err)
		}

		handle := handler.NewHandler()
		handle.AddFilters(filters...)
		handle.AddWriters(writers...)

		if cluster.Input.Kafka != nil {
			for j, k := range input.ToKqConf(cluster.Input.Kafka) {
				queue, err := kq.NewQueue(k, handle)
				if err != nil {
					return fmt.Errorf("initialize cluster %d kafka queue %d: %w", i, j, err)
				}
				initialized = append(initialized, queue)
				group.Add(queue)
			}
		}

		if cluster.Input.Syslogs != nil {
			for j, s := range cluster.Input.Syslogs {
				if s == nil {
					return fmt.Errorf("cluster %d syslog input %d is nil", i, j)
				}
				syslogService, err := input.NewSyslogService(s, handle)
				if err != nil {
					return fmt.Errorf("initialize cluster %d syslog input %d: %w", i, j, err)
				}
				initialized = append(initialized, syslogService)
				group.Add(syslogService)
			}
		}
	}

	return nil
}

func Run() error {
	if group == nil {
		return fmt.Errorf("service group not initialized")
	}

	group.Start()
	return nil
}

func Stop() {
	if group == nil {
		return
	}

	group.Stop()
}
