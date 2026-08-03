package service

import (
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/apperror"
	"context"
)

type TrafficService struct {
	dao TrafficDataAccess
}

type TrafficDataAccess interface {
	List(ctx context.Context, query model.TrafficQuery) ([]model.NetworkTraffic, error)
}

type TrafficListOptions struct {
	PageNo   int
	PageSize int
	OrderBy  string
}

func NewTrafficService(dao TrafficDataAccess) *TrafficService {
	return &TrafficService{dao: dao}
}

func (s *TrafficService) List(ctx context.Context, options TrafficListOptions) ([]model.NetworkTraffic, error) {
	if s.dao == nil {
		return nil, apperror.New(errorcode.SystemError)
	}

	return s.dao.List(ctx, model.TrafficQuery{
		PageNo:   options.PageNo,
		PageSize: options.PageSize,
		OrderBy:  options.OrderBy,
	})
}
