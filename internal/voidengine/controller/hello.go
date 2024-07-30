package controller

import (
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Ping struct {
}

func NewPing() *Ping {
	return &Ping{}
}

// @Summary Ping Example
// @Description Ping example endpoint
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (p *Ping) PingGet(c *gin.Context, e *env.Env) error {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Get ping")
	return nil
}

func (p *Ping) PingPost(c *gin.Context, e *env.Env) error {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Post ping")
	return nil
}
