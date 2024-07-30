package controller

import (
	"BlackHole/internal/voidengine/message"
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
// @Description Add a User
// @Tags User
// @Accept json
// @Produce json
// @param user body message.ListUserRequest true "list user param"
// @Success 200 {object} response.ApiResponse
// @Router /v1/user [get]
func (u *User) ListUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.ListUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParam
	}
	log.Info(request)

	return response.ApiSuccess
}

// AddUser
// @Description Add a User
// @Tags User
// @Accept json
// @Produce json
// @param user body message.AddUserRequest true "add user param"
// @Success 200 {object} response.ApiResponse
// @Router /v1//user [post]
func (u *User) AddUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.AddUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParam
	}
	log.Info(request)
	return response.ApiSuccess
}

// ModifyUser
// @Description Add a User
// @Tags User
// @Accept json
// @Produce json
// @param user body message.ModifyUserRequest true "modify user param"
// @Success 200 {object} response.ApiResponse
// @Router /v1//user [put]
func (u *User) ModifyUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.ModifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParam
	}
	log.Info(request)

	return response.ApiSuccess
}

// DeleteUser
// @Description Delete a User
// @Tags User
// @Accept json
// @Produce json
// @param user body message.DeleteUserRequest true "delete user param"
// @Success 200 {object} response.ApiResponse
// @Router /v1//user [delete]
func (u *User) DeleteUser(c *gin.Context, e *env.Env) *response.ApiResponse {
	var request message.DeleteUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.InvalidParam
	}
	log.Info(request)

	return response.ApiSuccess
}
