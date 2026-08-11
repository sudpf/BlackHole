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

func TestErrorParamsAreDeepCopied(t *testing.T) {
	params := Params{
		"user": map[string]any{
			"name": "alice",
			"tags": []string{"owner"},
			"roles": []any{
				"admin",
				map[string]any{"scope": "read"},
			},
		},
	}
	err := NewWithParams(100001, params)

	params["user"].(map[string]any)["name"] = "bob"
	params["user"].(map[string]any)["tags"].([]string)[0] = "guest"
	params["user"].(map[string]any)["roles"].([]any)[1].(map[string]any)["scope"] = "write"

	got := err.Params()
	user := got["user"].(map[string]any)
	tags := user["tags"].([]string)
	roles := user["roles"].([]any)
	if user["name"] != "alice" {
		t.Fatalf("stored name = %v, want alice", user["name"])
	}
	if tags[0] != "owner" {
		t.Fatalf("stored tag = %v, want owner", tags[0])
	}
	if roles[1].(map[string]any)["scope"] != "read" {
		t.Fatalf("stored scope = %v, want read", roles[1].(map[string]any)["scope"])
	}

	user["name"] = "charlie"
	tags[0] = "operator"
	roles[1].(map[string]any)["scope"] = "delete"

	got = err.Params()
	user = got["user"].(map[string]any)
	tags = user["tags"].([]string)
	roles = user["roles"].([]any)
	if user["name"] != "alice" {
		t.Fatalf("returned params mutated stored name = %v, want alice", user["name"])
	}
	if tags[0] != "owner" {
		t.Fatalf("returned params mutated stored tag = %v, want owner", tags[0])
	}
	if roles[1].(map[string]any)["scope"] != "read" {
		t.Fatalf("returned params mutated stored scope = %v, want read", roles[1].(map[string]any)["scope"])
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
