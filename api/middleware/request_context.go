package middleware

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
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

		traceID := requestctx.ResolveTraceID(c.GetHeader(requestctx.HeaderTraceID))
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
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
