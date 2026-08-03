package openapi

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewHTTPServerDoesNotMutateGinGlobals(t *testing.T) {
	originalMode := gin.Mode()
	originalWriter := gin.DefaultWriter
	originalErrorWriter := gin.DefaultErrorWriter
	t.Cleanup(func() {
		gin.SetMode(originalMode)
		gin.DefaultWriter = originalWriter
		gin.DefaultErrorWriter = originalErrorWriter
	})

	writer := &bytes.Buffer{}
	errorWriter := &bytes.Buffer{}
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = writer
	gin.DefaultErrorWriter = errorWriter

	if _, err := NewHTTPServer("127.0.0.1:8080", filepath.Join(t.TempDir(), "api.log"), "1m", time.Second); err != nil {
		t.Fatalf("NewHTTPServer error = %v", err)
	}

	if got := gin.Mode(); got != gin.TestMode {
		t.Fatalf("gin mode = %q, want %q", got, gin.TestMode)
	}
	if gin.DefaultWriter != writer {
		t.Fatal("gin DefaultWriter was mutated")
	}
	if gin.DefaultErrorWriter != errorWriter {
		t.Fatal("gin DefaultErrorWriter was mutated")
	}
}
