package errorcode

import (
	"BlackHole/pkg/apperror"
	"net/http"
)

const (
	Success                     = apperror.Success
	APINotFound   apperror.Code = 1
	InvalidParams apperror.Code = 2
	SystemError   apperror.Code = 3
	Unauthorized  apperror.Code = 4
	Forbidden     apperror.Code = 5

	InvalidUserName apperror.Code = 100001
	UserNotFound    apperror.Code = 100002
)

var Definitions = []apperror.Definition{
	{Code: Success, HTTPStatus: http.StatusOK, English: "Success", Chinese: "成功"},
	{Code: APINotFound, HTTPStatus: http.StatusNotFound, English: "API not found", Chinese: "API不存在"},
	{Code: InvalidParams, HTTPStatus: http.StatusBadRequest, English: "Invalid parameters", Chinese: "参数错误"},
	{Code: SystemError, HTTPStatus: http.StatusInternalServerError, English: "System error", Chinese: "系统错误"},
	{Code: Unauthorized, HTTPStatus: http.StatusUnauthorized, English: "Unauthorized", Chinese: "未认证"},
	{Code: Forbidden, HTTPStatus: http.StatusForbidden, English: "Forbidden", Chinese: "无权限"},
	{Code: InvalidUserName, HTTPStatus: http.StatusBadRequest, English: "Invalid username", Chinese: "用户名称不合法"},
	{Code: UserNotFound, HTTPStatus: http.StatusNotFound, English: "User {{.username}} does not exist", Chinese: "用户 {{.username}} 不存在"},
}
