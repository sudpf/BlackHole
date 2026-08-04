package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/api/voidengine/openapi/v1/message"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/apperror"

	"github.com/gin-gonic/gin"
)

// ListUser
// @Description List Users
// @Tags User
// @Accept json
// @Produce json
// @Param Accept-Language header string false "Language" default(zh)
// @Param pageNo query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(50)
// @Param orderBy query string false "排序方式[desc, asc]" default(desc)
// @Param username query string false "用户名"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [get]
func (h *Handler) ListUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request message.ListUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	users, err := h.services.User.List(ctx, service.UserListOptions{
		PageNo:   request.PageNo,
		PageSize: request.PageSize,
		OrderBy:  request.OrderBy,
		Username: request.UsernameFilter(),
	})
	if err != nil {
		return response.Result{}, err
	}

	return response.OK(users), nil
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
func (h *Handler) AddUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request message.AddUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Add(ctx, service.AddUserInput{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
	}); err != nil {
		return response.Result{}, err
	}

	return response.OK(nil), nil
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
func (h *Handler) ModifyUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request message.ModifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Modify(ctx, service.ModifyUserInput{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
	}); err != nil {
		return response.Result{}, err
	}

	return response.OK(nil), nil
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
func (h *Handler) DeleteUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request message.DeleteUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Delete(ctx, request.Username); err != nil {
		return response.Result{}, err
	}

	return response.OK(nil), nil
}
