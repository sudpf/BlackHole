package wrapper

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

type WrapperHandlerFunc func(*gin.Context, *env.Env)

func WrapperEnvFunc(handler WrapperHandlerFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = constant.LangEnglish
		}

		env := env.NewEnv(lang, c.ClientIP())

		handler(c, env)
	}
}
