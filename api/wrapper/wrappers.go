package wrapper

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*gin.Context) (response.Result, error)

func Adapt(catalog *apperror.Catalog, handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		requestEnv, ok := env.FromContext(c.Request.Context())
		if !ok {
			requestEnv = env.NewFromContext(c.Request.Context())
		}
		message, err := catalog.Localize(requestEnv, apperror.MessageID(apperror.Success), nil)
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
