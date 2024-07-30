package handler

import (
	"BlackHole/internal/voidengine/controller"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingGet(c *gin.Context, e *env.Env) {
	ping := controller.NewPing()
	res := ping.PingGet(c, e)
	c.JSON(http.StatusOK, res)
}

func PingPost(c *gin.Context, e *env.Env) {
	ping := controller.NewPing()
	res := ping.PingPost(c, e)
	c.JSON(http.StatusOK, res)
}
