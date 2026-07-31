package middleware

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/logger"
	"BlackHole/pkg/requestctx"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestContext(timeout time.Duration) gin.HandlerFunc {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return func(c *gin.Context) {
		language := c.GetHeader("Accept-Language")
		if language == "" {
			language = constant.LangEnglish
		}

		traceID, err := requestctx.ResolveTraceID(c.GetHeader(requestctx.HeaderTraceID))
		if err != nil {
			logger.FromContext(c.Request.Context()).WithError(err).Error("generate trace id")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), timeout)
		defer cancel()

		ctx = requestctx.WithScope(ctx, requestctx.Scope{
			TraceID:  traceID,
			Language: language,
			ClientIP: c.ClientIP(),
		})
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestctx.HeaderTraceID, traceID)

		c.Next()
	}
}
