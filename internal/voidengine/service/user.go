package service

import (
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/logger"
	"context"
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

type UserListOptions struct {
	PageNo   int
	PageSize int
	OrderBy  string
	Username *string
}

type AddUserInput struct {
	Username string
	Password string
	Email    string
	Phone    string
}

type ModifyUserInput struct {
	Username string
	Password *string
	Email    *string
	Phone    *string
}

func NewUserService(dao UserDataAccess) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) List(ctx context.Context, options UserListOptions) ([]model.User, error) {
	if s.dao == nil {
		return nil, apperror.New(errorcode.SystemError)
	}

	return s.dao.List(ctx, model.UserQuery{
		PageNo:   options.PageNo,
		PageSize: options.PageSize,
		OrderBy:  options.OrderBy,
		Username: options.Username,
	})
}

func (s *UserService) Add(ctx context.Context, input AddUserInput) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	return s.dao.Create(ctx, &model.User{
		Name:     input.Username,
		Password: input.Password,
		Email:    input.Email,
		Phone:    input.Phone,
	})
}

func (s *UserService) Modify(ctx context.Context, input ModifyUserInput) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	user, err := s.dao.FindByName(ctx, input.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return apperror.NewWithParams(errorcode.UserNotFound, apperror.Params{
			"username": input.Username,
		})
	}

	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}

	return s.dao.Update(ctx, input.Username, user)
}

func (s *UserService) Delete(ctx context.Context, username string) error {
	if s.dao == nil {
		return apperror.New(errorcode.SystemError)
	}

	user, err := s.dao.FindByName(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		logger.FromContext(ctx).
			WithField("username", username).
			Warn("delete user not found")
		return apperror.NewWithParams(errorcode.UserNotFound, apperror.Params{
			"username": username,
		})
	}

	return s.dao.DeleteByName(ctx, username)
}
