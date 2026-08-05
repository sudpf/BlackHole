package apperror

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrorPreservesCodeAndCause(t *testing.T) {
	cause := errors.New("database unavailable")
	err := WrapWithParams(100001, cause, Params{"name": "alice"})

	appErr, ok := As(err)
	if !ok {
		t.Fatal("As did not find application error")
	}
	if appErr.Code() != 100001 {
		t.Fatalf("Code = %d, want 100001", appErr.Code())
	}
	if !errors.Is(err, cause) {
		t.Fatal("wrapped cause is not preserved")
	}
	if got := appErr.Params()["name"]; got != "alice" {
		t.Fatalf("parameter name = %v, want alice", got)
	}
}

func TestNewCatalogRejectsDuplicateCode(t *testing.T) {
	_, err := NewCatalog(
		Definition{Code: Success, HTTPStatus: http.StatusOK, English: "Success", Chinese: "成功"},
		Definition{Code: Success, HTTPStatus: http.StatusCreated, English: "Created", Chinese: "创建"},
	)
	if err == nil {
		t.Fatal("NewCatalog expected duplicate code error")
	}
}

func TestCatalogMessageIDsAreSorted(t *testing.T) {
	catalog, err := NewCatalog(
		Definition{Code: 100001, HTTPStatus: http.StatusNotFound, English: "User not found", Chinese: "用户不存在"},
		Definition{Code: Success, HTTPStatus: http.StatusOK, English: "Success", Chinese: "成功"},
		Definition{Code: 2, HTTPStatus: http.StatusInternalServerError, English: "System error", Chinese: "系统错误"},
	)
	if err != nil {
		t.Fatalf("NewCatalog error = %v", err)
	}

	got := catalog.MessageIDs()
	want := []string{"error_0", "error_2", "error_100001"}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("MessageIDs[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}
