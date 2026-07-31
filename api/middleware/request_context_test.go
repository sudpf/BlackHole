package middleware

import (
	"BlackHole/pkg/requestctx"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRequestContextDetachesClientCancellation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestContext(time.Second))
	router.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := ctx.Err(); err != nil {
			t.Errorf("work context was canceled: %v", err)
		}
		if got := requestctx.TraceID(ctx); got != "trace-123" {
			t.Errorf("trace ID = %q, want trace-123", got)
		}
		if got := requestctx.Language(ctx); got != "zh" {
			t.Errorf("language = %q, want zh", got)
		}
		if got := c.Writer.Header().Get(requestctx.HeaderTraceID); got != "trace-123" {
			t.Errorf("response trace ID = %q, want trace-123", got)
		}
	})

	clientCtx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(clientCtx)
	request.Header.Set(requestctx.HeaderTraceID, "trace-123")
	request.Header.Set("Accept-Language", "zh")
	router.ServeHTTP(httptest.NewRecorder(), request)
}
