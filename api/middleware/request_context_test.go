package middleware

import (
	"BlackHole/pkg/env"
	"BlackHole/pkg/requestctx"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type requestContextTestKey struct{}

func TestRequestContextDetachesClientCancellation(t *testing.T) {
	traceID := "4bf92f3577b34da6a3ce929d0e0e4736"
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestContext(time.Second))
	router.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		if err := ctx.Err(); err != nil {
			t.Errorf("work context was canceled: %v", err)
		}
		if got := ctx.Value(requestContextTestKey{}); got != "upstream" {
			t.Errorf("upstream context value = %v, want upstream", got)
		}
		if got := requestctx.TraceID(ctx); got != traceID {
			t.Errorf("trace ID = %q, want %s", got, traceID)
		}
		if got := requestctx.Language(ctx); got != "zh" {
			t.Errorf("language = %q, want zh", got)
		}
		if got := c.Writer.Header().Get(requestctx.HeaderTraceID); got != traceID {
			t.Errorf("response trace ID = %q, want %s", got, traceID)
		}
		requestEnv, ok := env.FromContext(ctx)
		if !ok {
			t.Fatal("env is missing from context")
		}
		if requestEnv.RequestId != traceID {
			t.Errorf("env request ID = %q, want %s", requestEnv.RequestId, traceID)
		}
		if requestEnv.Lang != "zh" {
			t.Errorf("env language = %q, want zh", requestEnv.Lang)
		}
	})

	clientCtx, cancel := context.WithCancel(context.Background())
	clientCtx = context.WithValue(clientCtx, requestContextTestKey{}, "upstream")
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(clientCtx)
	request.Header.Set(requestctx.HeaderTraceID, traceID)
	request.Header.Set("Accept-Language", "zh")
	router.ServeHTTP(httptest.NewRecorder(), request)
}
