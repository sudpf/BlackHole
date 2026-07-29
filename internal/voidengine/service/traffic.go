package service

import (
	"BlackHole/internal/voidengine/model"
	"errors"
)

var ErrDataPlanDBUnavailable = errors.New("data plan database is unavailable")

type TrafficService struct {
	dao TrafficDataAccess
}

type TrafficDataAccess interface {
	List(query model.TrafficQuery) ([]model.NetworkTraffic, error)
}

type TrafficListOptions struct {
	PageNo   int
	PageSize int
	OrderBy  string
}

func NewTrafficService(dao TrafficDataAccess) *TrafficService {
	return &TrafficService{dao: dao}
}

func (s *TrafficService) List(options TrafficListOptions) ([]model.NetworkTraffic, error) {
	if s.dao == nil {
		return nil, ErrDataPlanDBUnavailable
	}

	return s.dao.List(model.TrafficQuery{
		PageNo:   options.PageNo,
		PageSize: options.PageSize,
		OrderBy:  options.OrderBy,
	})
}
