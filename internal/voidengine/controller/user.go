package controller

import (
	"BlackHole/internal/voidengine/message"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/response"
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// ListUser
// @Description List Users
// @Tags User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param user query message.ListUserRequest true "list user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [get]
func (u *User) ListUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.ListUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err))
	}
	log.Info(request)

	var users []model.User

	if _, err := model.ControlPlanDB().Query(&users, map[string]interface{}{}); err != nil {
		return response.SytemError
	}

	return response.ApiSuccess.WithData(users)
}

// AddUser
// @Description Add a User
// @Tags User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param user body message.AddUserRequest true "add user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [post]
func (u *User) AddUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.AddUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err))
	}
	log.Info(request)

	user := &model.User{
		Name:     request.Username,
		Password: request.Password,
	}

	if err := model.ControlPlanDB().Insert(user); err != nil {
		return response.SytemError
	}
	return response.ApiSuccess
}

// ModifyUser
// @Description Modify a User
// @Tags User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param user body message.ModifyUserRequest true "modify user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [put]
func (u *User) ModifyUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.ModifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err))
	}
	log.Info(request)

	user := &model.User{}

	if err := model.ControlPlanDB().Update(user, nil); err != nil {
		return response.SytemError
	}
	return response.ApiSuccess
}

// DeleteUser
// @Description Delete a User
// @Tags User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @param user body message.DeleteUserRequest true "delete user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [delete]
func (u *User) DeleteUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.DeleteUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err))
	}
	log.Info(request)

	user := &model.User{}

	if err := model.ControlPlanDB().Delete(user, nil); err != nil {
		return response.SytemError
	}
	return response.ApiSuccess
}
