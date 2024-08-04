package controller

import (
	"BlackHole/internal/voidengine/message"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/response"
	"BlackHole/pkg/env"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type NetworkTraffic struct {
}

func NewNetworkTraffic() *NetworkTraffic {
	return &NetworkTraffic{}
}

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
func (u *NetworkTraffic) ListNetworkTraffic(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.ListNetworkTrafficRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err))
	}
	log.Info(request)

	//TODO for test
	traffic := model.NetworkTraffic{
		ID:              1, // This should be managed by your application logic
		Timestamp:       time.Now(),
		SourceIP:        "192.168.1.1",
		DestinationIP:   "192.168.1.2",
		SourcePort:      12345,
		DestinationPort: 80,
		Protocol:        "TCP",
		BytesIn:         2048,
		BytesOut:        1024,
		PacketCount:     15,
		Description:     "Normal traffic",
	}

	model.DataPlanDB().Insert(&traffic)

	conditions := make(map[string]interface{})
	conditions["PageNo"] = request.ListQueryBase.PageNo
	conditions["PageSize"] = request.ListQueryBase.PageSize
	if len(request.ListQueryBase.OrderBy) > 0 {
		conditions["OrderBy"] = request.ListQueryBase.OrderBy
	}

	var traffics []model.NetworkTraffic
	if _, err := model.DataPlanDB().Query(&traffics, conditions); err != nil {
		log.Errorf("query db err:", err)
		return response.SytemError
	}

	return response.ApiSuccess.WithData(traffics)
}
