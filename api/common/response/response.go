package response

import (
	"BlackHole/pkg/apperror"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Details any    `json:"details,omitempty"`
}

type Result struct {
	Status int
	Data   any
}

func OK(data any) Result {
	return Result{Status: http.StatusOK, Data: data}
}

func Created(data any) Result {
	return Result{Status: http.StatusCreated, Data: data}
}

func WriteSuccess(c *gin.Context, message string, result Result) error {
	if result.Status < 200 || result.Status > 299 {
		return fmt.Errorf("invalid success HTTP status %d", result.Status)
	}

	c.JSON(result.Status, ApiResponse{
		Code:    int(apperror.Success),
		Message: message,
		Data:    result.Data,
	})
	return nil
}

func WriteError(c *gin.Context, status int, code apperror.Code, message string, details any) {
	c.JSON(status, ApiResponse{
		Code:    int(code),
		Message: message,
		Details: details,
	})
}
