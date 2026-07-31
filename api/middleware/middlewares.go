package middleware

import (
	"BlackHole/pkg/logger"
	"BlackHole/pkg/requestctx"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiLogMiddlewares(router *gin.Engine, logFile string, logSize string) error {
	output, err := logger.RotatingWriter(logFile, logSize)
	if err != nil {
		return err
	}

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: output, Formatter: func(param gin.LogFormatterParams) string {
		fields := map[string]interface{}{
			"time":       param.TimeStamp.Format("2006-01-02T15:04:05.000Z07:00"),
			"level":      "info",
			"trace_id":   requestctx.TraceID(param.Request.Context()),
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"protocol":   param.Request.Proto,
			"status":     param.StatusCode,
			"latency_ms": float64(param.Latency.Microseconds()) / 1000,
			"user_agent": param.Request.UserAgent(),
			"error":      param.ErrorMessage,
		}

		message, err := json.Marshal(fields)
		if err != nil {
			return fmt.Sprintf(`{"level":"error","error":%q}`+"\n", err.Error())
		}
		return string(message) + "\n"
	}}))

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.FromContext(c.Request.Context()).
			WithField("panic", fmt.Sprint(recovered)).
			Error("panic recovered")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	return nil
}
