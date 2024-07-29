package voidengine

import (
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// @Summary Ping Example
// @Description Ping example endpoint
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func PingGet(c *gin.Context, e *env.Env) {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Get ping")
	c.JSON(http.StatusOK, gin.H{
		"message": "get pong",
	})
}

func PingPost(c *gin.Context, e *env.Env) {
	log.WithFields(logrus.Fields{"clientip": e.ClientIp}).Error("Post ping")
	c.JSON(http.StatusOK, gin.H{
		"message": "post pong",
	})
}
