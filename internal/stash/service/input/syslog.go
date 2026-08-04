package input

import (
	"BlackHole/internal/stash/config"
	"BlackHole/internal/stash/service/handler"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/threading"
	"gopkg.in/mcuadros/go-syslog.v2"
)

type (
	ConsumeHandle func(key, value string) error

	ConsumeHandler interface {
		Consume(ctx context.Context, key, value string) error
	}

	SyslogService struct {
		ctx              context.Context
		conf             *config.SyslogServiceConf
		server           *syslog.Server
		channel          handler.LogPartsChannel
		handler          ConsumeHandler
		consumerRoutines *threading.RoutineGroup
		stopOnce         sync.Once
	}
)

func NewSyslogService(ctx context.Context, c *config.SyslogServiceConf, cHandler ConsumeHandler) (*SyslogService, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}
	channel := make(handler.LogPartsChannel, 10000)
	syslogHandler := handler.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(syslogHandler)

	switch c.Protocol {
	case "Udp":
		if err := server.ListenUDP(fmt.Sprintf("%s:%d", c.Address, c.Port)); err != nil {
			return nil, fmt.Errorf("create UDP listener: %w", err)
		}
	case "Tcp":
		if c.Ssl == "on" {
			if err := server.ListenTCPTLS(fmt.Sprintf("%s:%d", c.Address, c.Port), nil); err != nil {
				return nil, fmt.Errorf("create TCP TLS listener: %w", err)
			}
		} else {
			if err := server.ListenTCP(fmt.Sprintf("%s:%d", c.Address, c.Port)); err != nil {
				return nil, fmt.Errorf("create TCP listener: %w", err)
			}
		}
	case "Unixgram":
		os.Remove(c.Address)
		log.Infof("listen Unixgram on:%v", c.Address)

		if err := server.ListenUnixgram(c.Address); err != nil {
			return nil, fmt.Errorf("create unixgram listener: %w", err)
		}
	default:
		return nil, fmt.Errorf("unexpected syslog protocol: %s", c.Protocol)
	}

	if err := server.Boot(); err != nil {
		_ = server.Kill()
		return nil, fmt.Errorf("boot syslog server: %w", err)
	}

	return &SyslogService{
		ctx:              ctx,
		conf:             c,
		server:           server,
		channel:          channel,
		handler:          cHandler,
		consumerRoutines: threading.NewRoutineGroup(),
	}, nil
}

func (s *SyslogService) Start() {
	for i := 0; i < s.conf.Processors; i++ {
		s.consumerRoutines.Run(func() {
			log.Infof("Start syslog process [%d]", i+1)
			for logParts := range s.channel {
				select {
				case <-s.ctx.Done():
					log.Infof("Routine Run Over")
					return
				default:
				}
				logMap, err := jsoniter.MarshalToString(logParts)
				if err != nil {
					log.Warnf("Marshal error:%v parts[%v]", err, logParts)
					continue
				}

				if err := s.handler.Consume(s.ctx, strconv.FormatUint(1, 10), logMap); err != nil {
					log.Warnf("Consume log err:%v", err)
				}
			}
			log.Infof("Routine Run Over")
		})
	}
	s.consumerRoutines.Wait()
}

func (s *SyslogService) Stop() {
	s.stopOnce.Do(func() {
		log.Infof("Stop syslog service...")
		if s.server != nil {
			if err := s.server.Kill(); err != nil {
				log.Warnf("Kill syslog server error:%v", err)
			}
			s.server.Wait()
		}
		close(s.channel)
		if s.conf.Protocol == "Unixgram" {
			os.Remove(s.conf.Address)
		}
	})
}
