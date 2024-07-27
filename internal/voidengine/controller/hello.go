package controller

import (
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

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
