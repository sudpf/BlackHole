package service

import (
	"BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service/filter"
	"BlackHole/internal/stash/service/handler"
	"BlackHole/internal/stash/service/input"
	"BlackHole/internal/stash/service/output"
	"context"
	"fmt"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/proc"
	zeroservice "github.com/zeromicro/go-zero/core/service"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, cfg *config.Config) (_ *Stash, err error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	proc.SetTimeToForceQuit(cfg.GracePeriod)

	stash := &Stash{group: zeroservice.NewServiceGroup()}
	var initialized []zeroservice.Service
	defer func() {
		if err == nil {
			return
		}
		stopServices(initialized)
		_ = closeHandlers(ctx, stash.handlers)
	}()

	for i, cluster := range cfg.Clusters {
		clusterServices, handle, err := b.buildCluster(ctx, i, cluster)
		if err != nil {
			return nil, err
		}

		stash.handlers = append(stash.handlers, handle)
		for _, svc := range clusterServices {
			initialized = append(initialized, svc)
			stash.group.Add(svc)
		}
	}

	return stash, nil
}

func (b *Builder) buildCluster(ctx context.Context, index int, cluster *config.ClusterConf) ([]zeroservice.Service, *handler.MessageHandler, error) {
	if cluster == nil {
		return nil, nil, fmt.Errorf("cluster %d is nil", index)
	}
	if cluster.Input == nil {
		return nil, nil, fmt.Errorf("cluster %d input is nil", index)
	}

	handle, err := b.buildHandler(ctx, index, cluster)
	if err != nil {
		return nil, nil, err
	}

	services, err := b.buildInputServices(ctx, index, cluster.Input, handle)
	if err != nil {
		_ = handle.Close(ctx)
		return nil, nil, err
	}

	return services, handle, nil
}

func (b *Builder) buildHandler(ctx context.Context, clusterIndex int, cluster *config.ClusterConf) (*handler.MessageHandler, error) {
	writers, err := output.NewWriters(ctx, cluster.Output)
	if err != nil {
		return nil, fmt.Errorf("initialize cluster %d writers: %w", clusterIndex, err)
	}

	handle := handler.NewHandler()
	handle.AddFilters(filter.CreateFilters(cluster)...)
	handle.AddWriters(writers...)
	return handle, nil
}

func (b *Builder) buildInputServices(ctx context.Context, clusterIndex int, inputCfg *config.InputConf, handle *handler.MessageHandler) ([]zeroservice.Service, error) {
	var services []zeroservice.Service

	if inputCfg.Kafka != nil {
		for i, k := range input.ToKqConf(inputCfg.Kafka) {
			queue, err := kq.NewQueue(k, handle)
			if err != nil {
				stopServices(services)
				return nil, fmt.Errorf("initialize cluster %d kafka queue %d: %w", clusterIndex, i, err)
			}
			services = append(services, queue)
		}
	}

	for i, syslogCfg := range inputCfg.Syslogs {
		if syslogCfg == nil {
			stopServices(services)
			return nil, fmt.Errorf("cluster %d syslog input %d is nil", clusterIndex, i)
		}
		syslogService, err := input.NewSyslogService(ctx, syslogCfg, handle)
		if err != nil {
			stopServices(services)
			return nil, fmt.Errorf("initialize cluster %d syslog input %d: %w", clusterIndex, i, err)
		}
		services = append(services, syslogService)
	}

	return services, nil
}

func stopServices(services []zeroservice.Service) {
	for _, svc := range services {
		svc.Stop()
	}
}
