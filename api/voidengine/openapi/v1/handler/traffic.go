package handler

import (
	"BlackHole/internal/voidengine/controller"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListNetworkTraffic(c *gin.Context, e *env.Env) {
	traffic := controller.NewNetworkTraffic()
	res := traffic.ListNetworkTraffic(c, e)
	if res.Code != 0 {
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, res)
}
