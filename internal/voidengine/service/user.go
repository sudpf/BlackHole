package service

import (
	"BlackHole/internal/voidengine/model"
	"errors"
)

var (
	ErrControlPlanDBUnavailable = errors.New("control plan database is unavailable")
	ErrUserNotFound             = errors.New("user not found")
)

type UserService struct{}

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

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) List(options UserListOptions) ([]model.User, error) {
	database := model.ControlPlanDB()
	if database == nil {
		return nil, ErrControlPlanDBUnavailable
	}

	conditions := map[string]interface{}{
		"PageNo":   options.PageNo,
		"PageSize": options.PageSize,
	}
	if options.Username != nil {
		conditions["name"] = *options.Username
	}
	if options.OrderBy != "" {
		conditions["OrderBy"] = options.OrderBy
	}

	var users []model.User
	if _, err := database.Query(&users, conditions); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) Add(input AddUserInput) error {
	database := model.ControlPlanDB()
	if database == nil {
		return ErrControlPlanDBUnavailable
	}

	return database.Insert(&model.User{
		Name:     input.Username,
		Password: input.Password,
		Email:    input.Email,
		Phone:    input.Phone,
	})
}

func (s *UserService) Modify(input ModifyUserInput) error {
	database := model.ControlPlanDB()
	if database == nil {
		return ErrControlPlanDBUnavailable
	}

	var users []model.User
	if _, err := database.QueryEx(&users, &model.User{Name: input.Username}); err != nil {
		return err
	}
	if len(users) == 0 {
		return ErrUserNotFound
	}

	user := users[0]
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}

	return database.Update(&user, map[string]interface{}{"name": input.Username})
}

func (s *UserService) Delete(username string) error {
	database := model.ControlPlanDB()
	if database == nil {
		return ErrControlPlanDBUnavailable
	}

	return database.Delete(&model.User{}, map[string]interface{}{"name": username})
}
