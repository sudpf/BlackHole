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
	dao := model.GetTrafficDAO()
	if dao == nil {
		return nil, ErrDataPlanDBUnavailable
	}

	return dao.List(model.TrafficQuery{
		PageNo:   options.PageNo,
		PageSize: options.PageSize,
		OrderBy:  options.OrderBy,
	})
}
