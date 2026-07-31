package wrapper

import (
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

type WrapperHandlerFunc func(*gin.Context, *env.Env)

func WrapperEnvFunc(provider *env.Provider, handler WrapperHandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		handler(c, provider.NewEnvFromContext(c.Request.Context()))
	}
}
