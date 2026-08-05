package openapi_test

import (
	"BlackHole/api/common/response"
	"BlackHole/api/voidengine/openapi"
	v1handler "BlackHole/api/voidengine/openapi/v1/handler"
	v1router "BlackHole/api/voidengine/openapi/v1/router"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/service"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/auth"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

type authenticatorFunc func(context.Context, auth.Request) (*auth.Principal, error)

func (f authenticatorFunc) Authenticate(ctx context.Context, request auth.Request) (*auth.Principal, error) {
	return f(ctx, request)
}

func TestFrameworkResponses(t *testing.T) {
	apiServer, err := openapi.NewHTTPServer(
		"127.0.0.1:8080",
		filepath.Join(t.TempDir(), "api.log"),
		"1m",
		time.Second,
	)
	if err != nil {
		t.Fatalf("NewHTTPServer error = %v", err)
	}

	services := &service.Services{
		User:    service.NewUserService(nil),
		Traffic: service.NewTrafficService(nil),
	}
	v1router.RegisterRoutes(apiServer, v1handler.New(services))
	httpServer := apiServer.HTTPServer()

	tests := []struct {
		name        string
		path        string
		language    string
		wantStatus  int
		wantCode    int
		wantMessage string
		wantData    string
		wantDetails bool
	}{
		{
			name:        "ping success",
			path:        "/ping",
			language:    "zh-CN,zh;q=0.9",
			wantStatus:  http.StatusOK,
			wantCode:    0,
			wantMessage: "成功",
			wantData:    "Get ping",
		},
		{
			name:        "user validation error",
			path:        "/v1/user?pageNo=-1",
			language:    "zh",
			wantStatus:  http.StatusBadRequest,
			wantCode:    2,
			wantMessage: "参数错误",
			wantDetails: true,
		},
		{
			name:        "traffic validation error",
			path:        "/v1/traffic?pageSize=101",
			language:    "en",
			wantStatus:  http.StatusBadRequest,
			wantCode:    2,
			wantMessage: "Invalid parameters",
			wantDetails: true,
		},
		{
			name:        "unknown route",
			path:        "/missing",
			language:    "en",
			wantStatus:  http.StatusNotFound,
			wantCode:    1,
			wantMessage: "API not found",
		},
		{
			name:        "service error",
			path:        "/v1/user",
			language:    "en",
			wantStatus:  http.StatusInternalServerError,
			wantCode:    3,
			wantMessage: "System error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.Header.Set("Accept-Language", test.language)
			recorder := httptest.NewRecorder()
			httpServer.Handler.ServeHTTP(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}

			var body response.ApiResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if int(body.Code) != test.wantCode {
				t.Fatalf("code = %d, want %d", body.Code, test.wantCode)
			}
			if body.Message != test.wantMessage {
				t.Fatalf("message = %q, want %q", body.Message, test.wantMessage)
			}
			if test.wantData != "" && body.Data != test.wantData {
				t.Fatalf("data = %v, want %q", body.Data, test.wantData)
			}
			if test.wantDetails && body.Details == nil {
				t.Fatal("details are missing")
			}
		})
	}
}

func TestUserRoutesRequireAuthWhenAuthenticatorIsConfigured(t *testing.T) {
	apiServer, err := openapi.NewHTTPServer(
		"127.0.0.1:8080",
		filepath.Join(t.TempDir(), "api.log"),
		"1m",
		time.Second,
	)
	if err != nil {
		t.Fatalf("NewHTTPServer error = %v", err)
	}

	services := &service.Services{
		User:    service.NewUserService(nil),
		Traffic: service.NewTrafficService(nil),
	}
	authenticator := authenticatorFunc(func(context.Context, auth.Request) (*auth.Principal, error) {
		return nil, apperror.New(errorcode.Unauthorized)
	})
	v1router.RegisterRoutes(apiServer, v1handler.New(services), v1router.Options{
		Authenticator: authenticator,
	})
	httpServer := apiServer.HTTPServer()

	userRequest := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	userRequest.Header.Set("Accept-Language", "en")
	userRecorder := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(userRecorder, userRequest)
	if userRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("user status = %d, want %d; body = %s", userRecorder.Code, http.StatusUnauthorized, userRecorder.Body.String())
	}
	userBody := decodeAPIResponse(t, userRecorder)
	if userBody.Code != 4 || userBody.Message != "Unauthorized" {
		t.Fatalf("user response = %+v, want unauthorized", userBody)
	}

	pingRequest := httptest.NewRequest(http.MethodGet, "/ping", nil)
	pingRequest.Header.Set("Accept-Language", "en")
	pingRecorder := httptest.NewRecorder()
	httpServer.Handler.ServeHTTP(pingRecorder, pingRequest)
	if pingRecorder.Code != http.StatusOK {
		t.Fatalf("ping status = %d, want %d; body = %s", pingRecorder.Code, http.StatusOK, pingRecorder.Body.String())
	}
}

func decodeAPIResponse(t *testing.T, recorder *httptest.ResponseRecorder) response.ApiResponse {
	t.Helper()

	var body response.ApiResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}
