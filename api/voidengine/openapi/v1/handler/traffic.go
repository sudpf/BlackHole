package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/pkg/apperror"

	"github.com/gin-gonic/gin"
)

// ListNetworkTraffic
// @Description List NetworkTraffics
// @Tags NetworkTraffic
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param traffic query contract.ListNetworkTrafficRequest true "list traffic param"
// @Success 200 {object} response.ApiResponse{data=[]contract.ListNetworkTrafficResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /v1/traffic [get]
func (h *Handler) ListNetworkTraffic(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request contract.ListNetworkTrafficRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	traffics, err := h.services.Traffic.List(ctx, request)
	if err != nil {
		return response.Result{}, err
	}

	return response.OK(traffics), nil
}
