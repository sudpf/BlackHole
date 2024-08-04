package controller

import (
	"BlackHole/internal/voidengine/message"
	"BlackHole/internal/voidengine/model"
	"BlackHole/internal/voidengine/response"
	"BlackHole/pkg/common"
	"BlackHole/pkg/env"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"

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

	// 解析结构体
	reflectUser := &model.User{}
	dbSchema, err := schema.Parse(reflectUser, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return response.SytemError
	}

	reflectName := &reflectUser.Name
	field, ok := dbSchema.FieldsByName[common.FieldName(reflectUser, reflectName)]
	if !ok {
		return response.SytemError
	}

	conditions := make(map[string]interface{})
	if request.Username != nil {
		conditions[field.DBName] = request.Username
	}
	conditions["PageNo"] = request.ListQueryBase.PageNo
	conditions["PageSize"] = request.ListQueryBase.PageSize
	if len(request.ListQueryBase.OrderBy) > 0 {
		conditions["OrderBy"] = request.ListQueryBase.OrderBy
	}

	var users []model.User
	if _, err := model.ControlPlanDB().Query(&users, conditions); err != nil {
		log.Errorf("query db err:", err)
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
		Email:    request.Email,
		Phone:    request.Phone,
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

	var users []model.User
	if _, err := model.ControlPlanDB().QueryEx(&users, &model.User{Name: request.Username}); err != nil {
		return response.SytemError
	}
	if len(users) == 0 {
		return response.UserNotExist
	}

	modifyUser := users[0]
	if request.Password != nil {
		modifyUser.Password = *request.Password
	}
	if request.Email != nil {
		modifyUser.Email = *request.Email
	}
	if request.Phone != nil {
		modifyUser.Phone = *request.Phone
	}

	// 解析结构体
	reflectUser := &model.User{}
	dbSchema, err := schema.Parse(reflectUser, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return response.SytemError
	}

	reflectName := &reflectUser.Name
	field, ok := dbSchema.FieldsByName[common.FieldName(reflectUser, reflectName)]
	if !ok {
		return response.SytemError
	}

	conditions := make(map[string]interface{})
	conditions[field.DBName] = request.Username

	if err := model.ControlPlanDB().Update(&modifyUser, conditions); err != nil {
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

	// 解析结构体
	reflectUser := &model.User{}
	dbSchema, err := schema.Parse(reflectUser, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return response.SytemError
	}

	reflectName := &reflectUser.Name
	field, ok := dbSchema.FieldsByName[common.FieldName(reflectUser, reflectName)]
	if !ok {
		return response.SytemError
	}

	conditions := make(map[string]interface{})
	conditions[field.DBName] = request.Username

	if err := model.ControlPlanDB().Delete(&model.User{}, conditions); err != nil {
		return response.SytemError
	}
	return response.ApiSuccess
}
