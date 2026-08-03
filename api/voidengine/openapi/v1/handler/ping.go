package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PingGet
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /ping [get]
func (h *Handler) PingGet(c *gin.Context) (response.Result, error) {
	logger.FromContext(c.Request.Context()).Info("Get ping")
	return response.OK("Get ping"), nil
}

// PingPost
// @Description Ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /ping [post]
func (h *Handler) PingPost(c *gin.Context) (response.Result, error) {
	logger.FromContext(c.Request.Context()).Info("Post ping")
	return response.OK("Post ping"), nil
}
