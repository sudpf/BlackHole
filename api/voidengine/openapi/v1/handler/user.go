package handler

import (
	"BlackHole/api/common/response"
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/errorcode"
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
// @Success 200 {object} response.ApiResponse{data=[]contract.ListUserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [get]
func (h *Handler) ListUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request contract.ListUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	users, err := h.services.User.List(ctx, request)
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
// @param user body contract.AddUserRequest true "add user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [post]
func (h *Handler) AddUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request contract.AddUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Add(ctx, request); err != nil {
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
// @param user body contract.ModifyUserRequest true "modify user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [put]
func (h *Handler) ModifyUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request contract.ModifyUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Modify(ctx, request); err != nil {
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
// @param user body contract.DeleteUserRequest true "delete user param"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /v1/user [delete]
func (h *Handler) DeleteUser(c *gin.Context) (response.Result, error) {
	ctx := c.Request.Context()
	var request contract.DeleteUserRequest
	if err := c.ShouldBind(&request); err != nil {
		return response.Result{}, apperror.Wrap(errorcode.InvalidParams, err)
	}

	if err := h.services.User.Delete(ctx, request); err != nil {
		return response.Result{}, err
	}

	return response.OK(nil), nil
}
