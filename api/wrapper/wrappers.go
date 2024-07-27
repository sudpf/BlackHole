package wrapper

import (
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

type WrapperHandlerFunc func(*gin.Context, *env.Env)

func WrapperEnvFunc(handler WrapperHandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		env := &env.Env{}

		lang := c.Query("lang")
		if lang == "" {
			lang = c.GetHeader("Accept-Language")
		}

		env.Lang = lang
		env.ClientIp = c.ClientIP()
		handler(c, env)
	}
}
