package service

import (
	"BlackHole/internal/voidengine/model"
	"errors"
)

var ErrDataPlanDBUnavailable = errors.New("data plan database is unavailable")

type TrafficService struct{}

type TrafficListOptions struct {
	PageNo   int
	PageSize int
	OrderBy  string
}

func NewTrafficService() *TrafficService {
	return &TrafficService{}
}

func (s *TrafficService) List(options TrafficListOptions) ([]model.NetworkTraffic, error) {
	database := model.DataPlanDB()
	if database == nil {
		return nil, ErrDataPlanDBUnavailable
	}

	conditions := map[string]interface{}{
		"PageNo":   options.PageNo,
		"PageSize": options.PageSize,
	}
	if options.OrderBy != "" {
		conditions["OrderBy"] = options.OrderBy
	}

	var traffics []model.NetworkTraffic
	if _, err := database.Query(&traffics, conditions); err != nil {
		return nil, err
	}

	return traffics, nil
}
