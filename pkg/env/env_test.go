package env

import (
	"BlackHole/pkg/auth"
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
	"testing"
)

func TestNewFromContext(t *testing.T) {
	ctx := requestctx.WithScope(context.Background(), requestctx.Scope{
		TraceID:  "4bf92f3577b34da6a3ce929d0e0e4736",
		Language: "zh-CN,zh;q=0.9,en;q=0.8",
		ClientIP: "127.0.0.1",
	})

	ev := NewFromContext(ctx)
	if ev.Lang != "zh-CN,zh;q=0.9,en;q=0.8" {
		t.Fatalf("Lang = %q, want raw accept language", ev.Lang)
	}
	if ev.RequestId != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("RequestId = %q, want 4bf92f3577b34da6a3ce929d0e0e4736", ev.RequestId)
	}
	if ev.ClientIp != "127.0.0.1" {
		t.Fatalf("ClientIp = %q, want 127.0.0.1", ev.ClientIp)
	}
}

func TestEnvContextAndWithPrincipal(t *testing.T) {
	ev := &Env{Lang: constant.LangEnglish}
	ctx := WithContext(context.Background(), ev)

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("env is missing from context")
	}
	if got != ev {
		t.Fatal("env from context should be the stored value")
	}

	principal := &auth.Principal{Subject: "alice", Method: auth.MethodAPIKey}
	updated := ev.WithPrincipal(principal)
	if updated == ev {
		t.Fatal("WithPrincipal should return a copied env")
	}
	if updated.Principal != principal {
		t.Fatalf("principal = %+v, want %+v", updated.Principal, principal)
	}
	if ev.Principal != nil {
		t.Fatal("original env should not be modified")
	}
}
