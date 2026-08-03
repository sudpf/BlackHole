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
		Definition{Code: Success, HTTPStatus: http.StatusOK},
		Definition{Code: Success, HTTPStatus: http.StatusCreated},
	)
	if err == nil {
		t.Fatal("NewCatalog expected duplicate code error")
	}
}

func TestCatalogMessageIDsAreSorted(t *testing.T) {
	catalog, err := NewCatalog(
		Definition{Code: 100001, HTTPStatus: http.StatusNotFound},
		Definition{Code: Success, HTTPStatus: http.StatusOK},
		Definition{Code: 2, HTTPStatus: http.StatusBadRequest},
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
