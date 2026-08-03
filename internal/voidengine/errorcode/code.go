package errorcode

import "BlackHole/pkg/apperror"

const (
	APINotFound   apperror.Code = 1
	InvalidParams apperror.Code = 2
	SystemError   apperror.Code = 3

	InvalidUserName apperror.Code = 100001
	UserNotFound    apperror.Code = 100002
)
