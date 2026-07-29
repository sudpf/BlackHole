package service

import "BlackHole/internal/voidengine/model"

type Services struct {
	User    *UserService
	Traffic *TrafficService
}

func New(models *model.Models) *Services {
	var userDAO UserDataAccess
	if models.User != nil {
		userDAO = models.User
	}

	var trafficDAO TrafficDataAccess
	if models.Traffic != nil {
		trafficDAO = models.Traffic
	}

	return &Services{
		User:    NewUserService(userDAO),
		Traffic: NewTrafficService(trafficDAO),
	}
}
