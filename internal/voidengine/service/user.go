package service

import (
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/logger"
	"context"
	"strings"
)

type UserService struct {
	dao UserDataAccess
}

type UserDataAccess interface {
	List(ctx context.Context, query model.UserQuery) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	FindByName(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, username string, user *model.User) error
	DeleteByName(ctx context.Context, username string) error
}

func NewUserService(dao UserDataAccess) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) List(ctx context.Context, request contract.ListUserRequest) ([]contract.ListUserResponse, error) {
	if s.dao == nil {
		return nil, apperror.New(errorcode.SystemError)
	}

	var username *string
	if request.Username != nil && strings.TrimSpace(*request.Username) != "" {
		username = request.Username
	}

	users, err := s.dao.List(ctx, model.UserQuery{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	response := make([]contract.ListUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, contract.ListUserResponse{
			Username: user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
		})
	}
	return response, nil
}

func (s *UserService) Add(ctx context.Context, request contract.AddUserRequest) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return apperror.Wrap(errorcode.SystemError, err)
	}

	return s.dao.Create(ctx, &model.User{
		Name:     request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		Phone:    request.Phone,
	})
}

func (s *UserService) Modify(ctx context.Context, request contract.ModifyUserRequest) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	user, err := s.dao.FindByName(ctx, request.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.NewWithParams(errorcode.UserNotFound, apperror.Params{
			"username": request.Username,
		})
	}

	if request.Password != nil {
		hashedPassword, err := hashPassword(*request.Password)
		if err != nil {
			return apperror.Wrap(errorcode.SystemError, err)
		}
		user.Password = hashedPassword
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Phone != nil {
		user.Phone = *request.Phone
	}

	return s.dao.Update(ctx, request.Username, user)
}

func (s *UserService) Delete(ctx context.Context, request contract.DeleteUserRequest) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	user, err := s.dao.FindByName(ctx, request.Username)
	if err != nil {
		return err
	}
	if user == nil {
		logger.FromContext(ctx).
			WithField("username", request.Username).
			Warn("delete user not found")
		return apperror.NewWithParams(errorcode.UserNotFound, apperror.Params{
			"username": request.Username,
		})
	}

	return s.dao.DeleteByName(ctx, request.Username)
}
