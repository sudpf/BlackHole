package middleware

import (
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/logger"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err := fmt.Errorf("panic: %v", recovered)
				logger.FromContext(c.Request.Context()).
					WithError(err).
					WithField("stack", string(debug.Stack())).
					Error("panic recovered")

				_ = c.Error(apperror.Wrap(errorcode.SystemError, err))
				c.Abort()
			}
		}()

		c.Next()
	}
}
