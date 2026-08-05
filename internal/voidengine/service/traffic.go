package service

import (
	"BlackHole/internal/voidengine/contract"
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

func NewTrafficService(dao TrafficDataAccess) *TrafficService {
	return &TrafficService{dao: dao}
}

func (s *TrafficService) List(ctx context.Context, request contract.ListNetworkTrafficRequest) ([]contract.ListNetworkTrafficResponse, error) {
	if s.dao == nil {
		return nil, apperror.New(errorcode.SystemError)
	}

	traffics, err := s.dao.List(ctx, model.TrafficQuery{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
	})
	if err != nil {
		return nil, err
	}

	response := make([]contract.ListNetworkTrafficResponse, 0, len(traffics))
	for _, traffic := range traffics {
		response = append(response, contract.ListNetworkTrafficResponse{
			Timestamp:       traffic.Timestamp,
			SourceIP:        traffic.SourceIP,
			DestinationIP:   traffic.DestinationIP,
			SourcePort:      traffic.SourcePort,
			DestinationPort: traffic.DestinationPort,
			Protocol:        traffic.Protocol,
			BytesIn:         traffic.BytesIn,
			BytesOut:        traffic.BytesOut,
			PacketCount:     traffic.PacketCount,
			Description:     traffic.Description,
		})
	}
	return response, nil
}
