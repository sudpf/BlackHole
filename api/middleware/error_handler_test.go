package middleware

import (
	"BlackHole/api/common/response"
	"BlackHole/api/validation"
	"BlackHole/pkg/apperror"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	testInvalidParams apperror.Code = 2
	testSystemError   apperror.Code = 3
	testUserNotFound  apperror.Code = 100002
)

func TestErrorHandlerLocalizesApplicationError(t *testing.T) {
	router := newErrorTestRouter(t)
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(apperror.NewWithParams(testUserNotFound, apperror.Params{
			"username": "missing",
		}))
		c.Abort()
	})

	recorder := performErrorTestRequest(router, "zh-CN,zh;q=0.9")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	body := decodeErrorResponse(t, recorder)
	if body.Code != int(testUserNotFound) {
		t.Fatalf("code = %d, want %d", body.Code, testUserNotFound)
	}
	if body.Message != "用户 missing 不存在" {
		t.Fatalf("message = %q, want 用户 missing 不存在", body.Message)
	}
}

func TestErrorHandlerTranslatesValidationDetails(t *testing.T) {
	router := newErrorTestRouter(t)
	router.GET("/", func(c *gin.Context) {
		var input struct {
			Username string `form:"username" json:"username" binding:"required"`
		}
		if err := c.ShouldBindQuery(&input); err != nil {
			_ = c.Error(apperror.Wrap(testInvalidParams, err))
			c.Abort()
		}
	})

	recorder := performErrorTestRequest(router, "zh")
	body := decodeErrorResponse(t, recorder)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	details, ok := body.Details.(map[string]any)
	if !ok {
		t.Fatalf("details type = %T, want map[string]any", body.Details)
	}
	if _, exists := details["username"]; !exists {
		t.Fatalf("details = %v, want username validation error", details)
	}
}

func TestErrorHandlerConvertsUnknownError(t *testing.T) {
	router := newErrorTestRouter(t)
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(errors.New("database password must stay private"))
		c.Abort()
	})

	recorder := performErrorTestRequest(router, "en")
	body := decodeErrorResponse(t, recorder)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if body.Code != int(testSystemError) || body.Message != "System error" {
		t.Fatalf("response = %+v, want system error", body)
	}
}

func TestErrorHandlerRecoversPanic(t *testing.T) {
	router := newErrorTestRouter(t)
	router.GET("/", func(*gin.Context) {
		panic("unexpected")
	})

	recorder := performErrorTestRequest(router, "en")
	body := decodeErrorResponse(t, recorder)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if body.Code != int(testSystemError) {
		t.Fatalf("code = %d, want %d", body.Code, testSystemError)
	}
}

type testErrorRegistry struct {
	catalog    *apperror.Catalog
	systemCode apperror.Code
}

func (r testErrorRegistry) Catalog() *apperror.Catalog {
	return r.catalog
}

func (r testErrorRegistry) SystemErrorCode() apperror.Code {
	return r.systemCode
}

func newErrorTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	catalog, err := apperror.NewCatalog(
		apperror.Definition{Code: apperror.Success, HTTPStatus: http.StatusOK, English: "Success", Chinese: "成功"},
		apperror.Definition{Code: testInvalidParams, HTTPStatus: http.StatusBadRequest, English: "Invalid parameters", Chinese: "参数错误"},
		apperror.Definition{Code: testSystemError, HTTPStatus: http.StatusInternalServerError, English: "System error", Chinese: "系统错误"},
		apperror.Definition{Code: testUserNotFound, HTTPStatus: http.StatusNotFound, English: "User {{.username}} does not exist", Chinese: "用户 {{.username}} 不存在"},
	)
	if err != nil {
		t.Fatalf("NewCatalog error = %v", err)
	}
	registry := testErrorRegistry{catalog: catalog, systemCode: testSystemError}
	translator, err := validation.NewTranslator()
	if err != nil {
		t.Fatalf("NewTranslator error = %v", err)
	}
	router := gin.New()
	router.Use(ErrorHandler(registry, translator))
	router.Use(Recovery(registry))
	router.Use(RequestContext(time.Second))
	return router
}

func performErrorTestRequest(router *gin.Engine, language string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Language", language)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func decodeErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) response.ApiResponse {
	t.Helper()

	var body response.ApiResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}
