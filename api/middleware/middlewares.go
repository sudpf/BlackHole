package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitMiddlewares(router *gin.Engine) {
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: log.StandardLogger().Out, Formatter: func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("time=\"%s\" client=\"%s\" method=\"%s\" path=\"%s\" protocol=\"%s\" code=\"%d\" latency=\"%s\" useragent=\"%s\" error=\"%s\"\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}}))

	router.Use(gin.Recovery())
}
