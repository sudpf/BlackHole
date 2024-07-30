package handler

import (
	"BlackHole/internal/voidengine/controller"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingGet(c *gin.Context, e *env.Env) {
	ping := controller.NewPing()
	ping.PingGet(c, e)
	c.JSON(http.StatusOK, gin.H{
		"message": "get pong",
	})
}

func PingPost(c *gin.Context, e *env.Env) {
	ping := controller.NewPing()
	ping.PingPost(c, e)
	c.JSON(http.StatusOK, gin.H{
		"message": "post pong",
	})
}
