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

func Init(cfg *config.StashConfig) error {
	proc.SetTimeToForceQuit(cfg.GracePeriod)
	group = service.NewServiceGroup()

	for i, cluster := range cfg.Clusters {
		filters := filter.CreateFilters(cluster)

		writers, err := output.NewWriters(cluster.Output)
		if err != nil {
			return fmt.Errorf("initialize cluster %d writers: %w", i, err)
		}

		handle := handler.NewHandler()
		handle.AddFilters(filters...)
		handle.AddWriters(writers...)

		if cluster.Input.Kafka != nil {
			for _, k := range input.ToKqConf(cluster.Input.Kafka) {
				group.Add(kq.MustNewQueue(k, handle))
			}
		}

		if cluster.Input.Syslogs != nil {
			for _, s := range cluster.Input.Syslogs {
				group.Add(input.NewSyslogService(s, handle))
			}
		}
	}

	return nil
}

func Run() {
	group.Start()
}

func Stop() {
	group.Stop()
}
