package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/api/voidengine/openapi/v1/message"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/env"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
func (h *Handler) ListUser(c *gin.Context, e *env.Env) {
	var request message.ListUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	users, err := h.services.User.List(service.UserListOptions{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
		Username: request.Username,
	})
	if err != nil {
		respondUserServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess.WithData(users))
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
func (h *Handler) AddUser(c *gin.Context, e *env.Env) {
	var request message.AddUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	if err := h.services.User.Add(service.AddUserInput{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
	}); err != nil {
		respondUserServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess)
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
func (h *Handler) ModifyUser(c *gin.Context, e *env.Env) {
	var request message.ModifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	if err := h.services.User.Modify(service.ModifyUserInput{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
	}); err != nil {
		respondUserServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess)
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
func (h *Handler) DeleteUser(c *gin.Context, e *env.Env) {
	var request message.DeleteUserRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.InvalidParams.Tr(e).WithData(e.TranslatErrors(err)))
		return
	}

	if err := h.services.User.Delete(request.Username); err != nil {
		respondUserServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiSuccess)
}

func respondUserServiceError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrUserNotFound) {
		c.JSON(http.StatusBadRequest, response.UserNotExist)
		return
	}

	log.WithError(err).Error("Handle user request")
	c.JSON(http.StatusInternalServerError, response.SytemError)
}
