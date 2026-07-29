package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// PingGet
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /ping [get]
func (h *Handler) PingGet(c *gin.Context, e *env.Env) {
	log.WithField("clientip", e.ClientIp).Info("Get ping")
	c.JSON(http.StatusOK, response.ApiSuccess.Tr(e).WithData("Get ping"))
}

// PingPost
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /ping [post]
func (h *Handler) PingPost(c *gin.Context, e *env.Env) {
	log.WithField("clientip", e.ClientIp).Info("Post ping")
	c.JSON(http.StatusOK, response.ApiSuccess.Tr(e).WithData("Post ping"))
}
