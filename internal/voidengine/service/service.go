package service

import (
	"BlackHole/internal/voidengine/model"
	"errors"
)

type ErrorCode string

const (
	ErrorCodeDependencyUnavailable ErrorCode = "DEPENDENCY_UNAVAILABLE"
	ErrorCodeUserNotFound          ErrorCode = "USER_NOT_FOUND"
)

type Error struct {
	Code ErrorCode
}

func (e *Error) Error() string {
	return string(e.Code)
}

func NewError(code ErrorCode) error {
	return &Error{Code: code}
}

func IsErrorCode(err error, code ErrorCode) bool {
	var serviceError *Error
	return errors.As(err, &serviceError) && serviceError.Code == code
}

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
