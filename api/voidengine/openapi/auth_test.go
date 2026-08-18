package openapi

import (
	"BlackHole/api/middleware"
	"BlackHole/api/router"
	"BlackHole/api/validation"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/auth"
	"BlackHole/pkg/env"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type testAuthenticator func(context.Context, auth.Request) (*auth.Principal, error)

func (f testAuthenticator) Authenticate(ctx context.Context, request auth.Request) (*auth.Principal, error) {
	return f(ctx, request)
}

func TestRequireAuthWritesPrincipalToEnv(t *testing.T) {
	engine := newAuthTestEngine(t)
	principal := &auth.Principal{Subject: "alice", Method: auth.MethodAPIKey}
	authenticator := testAuthenticator(func(context.Context, auth.Request) (*auth.Principal, error) {
		return principal, nil
	})

	route := router.NewGetRoute("/protected", func(c *gin.Context) {
		requestEnv, ok := env.FromContext(c.Request.Context())
		if !ok {
			t.Fatal("env is missing from context")
		}
		if requestEnv.Principal != principal {
			t.Fatalf("principal = %+v, want %+v", requestEnv.Principal, principal)
		}
		c.Status(http.StatusNoContent)
	}, RequireAuth(authenticator))
	engine.GET(route.Path(), route.Handler())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusNoContent, recorder.Body.String())
	}
}

func TestRequireAuthMapsMissingCredentialToUnauthorized(t *testing.T) {
	engine := newAuthTestEngine(t)
	authenticator := testAuthenticator(func(context.Context, auth.Request) (*auth.Principal, error) {
		return nil, apperror.New(errorcode.Unauthorized)
	})

	route := router.NewGetRoute("/protected", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	}, RequireAuth(authenticator))
	engine.GET(route.Path(), route.Handler())

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Accept-Language", "en")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusUnauthorized, recorder.Body.String())
	}
}

func newAuthTestEngine(t *testing.T) *gin.Engine {
	t.Helper()

	catalog, err := apperror.NewCatalog(errorcode.Definitions...)
	if err != nil {
		t.Fatalf("error catalog initialization error = %v", err)
	}
	translator, err := validation.NewTranslator()
	if err != nil {
		t.Fatalf("validation translator initialization error = %v", err)
	}

	engine := gin.New()
	engine.Use(middleware.ErrorHandler(catalog, translator))
	engine.Use(middleware.RequestContext(time.Second))
	engine.NoRoute(func(c *gin.Context) {
		_ = c.Error(apperror.New(errorcode.APINotFound))
		c.Abort()
	})
	return engine
}
