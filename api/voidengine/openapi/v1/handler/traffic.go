package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/api/voidengine/openapi/v1/message"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/env"
	"BlackHole/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListNetworkTraffic
// @Description List NetworkTraffics
// @Tags NetworkTraffic
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param traffic query message.ListNetworkTrafficRequest true "list traffic param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/traffic [get]
func (h *Handler) ListNetworkTraffic(c *gin.Context, e *env.Env) {
	ctx := c.Request.Context()
	var request message.ListNetworkTrafficRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	traffics, err := h.services.Traffic.List(ctx, service.TrafficListOptions{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
	})
	if err != nil {
		logger.FromContext(ctx).WithError(err).Error("List network traffic")
		c.JSON(http.StatusInternalServerError, response.SystemError.Tr(e))
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess.Tr(e).WithData(traffics))
}
