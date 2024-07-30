package controller

import (
	"BlackHole/internal/voidengine/response"
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

// Ping
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /v1/ping [get]
func (p *Ping) PingGet(c *gin.Context, e *env.Env) *response.ApiResponse {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Get ping")
	return response.ApiSuccess.WithData("Get ping")
}

func (p *Ping) PingPost(c *gin.Context, e *env.Env) *response.ApiResponse {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Post ping")
	return response.ApiSuccess.WithData("Post ping")
}
