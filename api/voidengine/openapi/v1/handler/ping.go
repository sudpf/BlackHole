package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/env"
	"BlackHole/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingGet
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /ping [get]
func (h *Handler) PingGet(c *gin.Context, e *env.Env) {
	logger.FromContext(c.Request.Context()).Info("Get ping")
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
	logger.FromContext(c.Request.Context()).Info("Post ping")
	c.JSON(http.StatusOK, response.ApiSuccess.Tr(e).WithData("Post ping"))
}
