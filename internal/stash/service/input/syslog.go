package input

import (
	"BlackHole/internal/stash/service/handler"
	"BlackHole/pkg/config"
	"context"
	"fmt"
	"os"
	"strconv"

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
		conf             *config.SyslogServiceConf
		channel          handler.LogPartsChannel
		handler          ConsumeHandler
		consumerRoutines *threading.RoutineGroup
	}
)

func NewSyslogService(c *config.SyslogServiceConf, cHandler ConsumeHandler) SyslogService {
	channel := make(handler.LogPartsChannel, 10000)
	syslogHandler := handler.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(syslogHandler)

	switch c.Protocol {
	case "Udp":
		if err := server.ListenUDP(fmt.Sprintf("%s:%d", c.Address, c.Port)); err != nil {
			log.Errorf("Create UDP listen error:%v", err)
		}
	case "Tcp":
		if c.Ssl == "on" {
			if err := server.ListenTCPTLS(fmt.Sprintf("%s:%d", c.Address, c.Port), nil); err != nil {
				log.Errorf("Create TCPTLS listen error:%v", err)
			}
		} else {
			if err := server.ListenTCP(fmt.Sprintf("%s:%d", c.Address, c.Port)); err != nil {
				log.Errorf("Create TCP listen error:%v", err)
			}
		}
	case "Unixgram":
		os.Remove(c.Address)
		log.Infof("listen Unixgram on:%v", c.Address)

		if err := server.ListenUnixgram(c.Address); err != nil {
			log.Errorf("create unixgram listen error:%v", err)
		}
	default:
		log.Errorf("Unexpect syslog protocol:%v", c.Protocol)
	}

	server.Boot()

	return SyslogService{
		conf:             c,
		channel:          channel,
		handler:          cHandler,
		consumerRoutines: threading.NewRoutineGroup(),
	}
}

func (s SyslogService) Start() {
	for i := 0; i < s.conf.Processors; i++ {
		s.consumerRoutines.Run(func() {
			log.Infof("Start syslog process [%d]", i+1)
			ctx := context.TODO()
			for logParts := range s.channel {
				logMap, err := jsoniter.MarshalToString(logParts)
				if err != nil {
					log.Warn("Marshal error:%v parts[%v]", err, logParts)
					continue
				}

				if err := s.handler.Consume(ctx, strconv.FormatUint(1, 10), logMap); err != nil {
					log.Warn("Consume log err:%v", err)
				}
			}
			log.Infof("Routine Run Over")
		})
	}
	s.consumerRoutines.Wait()
}

func (s SyslogService) Stop() {
	log.Infof("Stop syslog service...")
	if s.conf.Protocol == "Unixgram" {
		os.Remove(s.conf.Address)
	}
}
