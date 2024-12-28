package service

import (
	"BlackHole/internal/stash/service/filter"
	"BlackHole/internal/stash/service/handler"
	"BlackHole/internal/stash/service/input"
	"BlackHole/internal/stash/service/output"
	"BlackHole/pkg/config"

	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
)

var (
	group *service.ServiceGroup
)

func Init() {
	proc.SetTimeToForceQuit(config.GlobalStashConfig.GracePeriod)
	group = service.NewServiceGroup()

	for _, cluster := range config.GlobalStashConfig.Clusters {
		filters := filter.CreateFilters(cluster)

		writers, err := output.NewWriters(cluster.Output)
		if err != nil {
			log.Warn("NewWriters err: %v", err)
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
}

func Run() {
	group.Start()
}

func Stop() {
	group.Stop()
}
