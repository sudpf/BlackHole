package service

import (
	"BlackHole/internal/voidengine/model"
	"errors"
)

var (
	ErrControlPlanDBUnavailable = errors.New("control plan database is unavailable")
	ErrUserNotFound             = errors.New("user not found")
)

type UserService struct {
	dao UserDataAccess
}

type UserDataAccess interface {
	List(query model.UserQuery) ([]model.User, error)
	Create(user *model.User) error
	FindByName(username string) (*model.User, error)
	Update(username string, user *model.User) error
	DeleteByName(username string) error
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

func (s *UserService) List(options UserListOptions) ([]model.User, error) {
	if s.dao == nil {
		return nil, ErrControlPlanDBUnavailable
	}

	return s.dao.List(model.UserQuery{
		PageNo:   options.PageNo,
		PageSize: options.PageSize,
		OrderBy:  options.OrderBy,
		Username: options.Username,
	})
}

func (s *UserService) Add(input AddUserInput) error {
	if s.dao == nil {
		return ErrControlPlanDBUnavailable
	}

	return s.dao.Create(&model.User{
		Name:     input.Username,
		Password: input.Password,
		Email:    input.Email,
		Phone:    input.Phone,
	})
}

func (s *UserService) Modify(input ModifyUserInput) error {
	if s.dao == nil {
		return ErrControlPlanDBUnavailable
	}

	user, err := s.dao.FindByName(input.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
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

	return s.dao.Update(input.Username, user)
}

func (s *UserService) Delete(username string) error {
	if s.dao == nil {
		return ErrControlPlanDBUnavailable
	}

	return s.dao.DeleteByName(username)
}
