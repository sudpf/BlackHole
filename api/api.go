package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitApi() {
	// Creates a router without any middleware by default
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	//router.Use(gin.LoggerWithWriter(log.StandardLogger().Out))

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

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
