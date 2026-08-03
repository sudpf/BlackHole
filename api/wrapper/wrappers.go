package wrapper

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*gin.Context) (response.Result, error)

func Adapt(provider *env.Provider, handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		requestEnv := provider.NewEnvFromContext(c.Request.Context())
		message, err := requestEnv.Localize(apperror.MessageID(apperror.Success), nil)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if err := response.WriteSuccess(c, message, result); err != nil {
			_ = c.Error(err)
			c.Abort()
		}
	}
}
