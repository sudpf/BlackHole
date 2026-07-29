package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/api/voidengine/openapi/v1/message"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/env"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	var request message.ListNetworkTrafficRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	traffics, err := h.services.Traffic.List(service.TrafficListOptions{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
	})
	if err != nil {
		log.WithError(err).Error("List network traffic")
		c.JSON(http.StatusInternalServerError, response.SytemError)
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess.WithData(traffics))
}
