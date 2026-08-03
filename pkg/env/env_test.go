package env

import (
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
	"testing"
)

func TestProviderNewEnvFromContext(t *testing.T) {
	provider, err := NewProvider(
		map[string]string{"hello": "hello"},
		map[string]string{"hello": "你好"},
	)
	if err != nil {
		t.Fatalf("NewProvider error = %v", err)
	}

	ctx := requestctx.WithScope(context.Background(), requestctx.Scope{
		TraceID:  "4bf92f3577b34da6a3ce929d0e0e4736",
		Language: "zh-CN,zh;q=0.9,en;q=0.8",
		ClientIP: "127.0.0.1",
	})

	ev := provider.NewEnvFromContext(ctx)
	if ev.Lang != constant.LangChinese {
		t.Fatalf("Lang = %q, want %q", ev.Lang, constant.LangChinese)
	}
	if ev.RequestId != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("RequestId = %q, want 4bf92f3577b34da6a3ce929d0e0e4736", ev.RequestId)
	}
	if ev.ClientIp != "127.0.0.1" {
		t.Fatalf("ClientIp = %q, want 127.0.0.1", ev.ClientIp)
	}
	got, err := ev.Localize("hello", nil)
	if err != nil {
		t.Fatalf("Localize error = %v", err)
	}
	if got != "你好" {
		t.Fatalf("Localize = %q, want 你好", got)
	}
}

func TestNewProviderRejectsIncompleteMessages(t *testing.T) {
	_, err := NewProvider(
		map[string]string{"hello": "hello"},
		map[string]string{},
	)
	if err == nil {
		t.Fatal("NewProvider expected missing translation error")
	}
}

func TestInitValidatorTranslationsRejectsNilProvider(t *testing.T) {
	if err := InitValidatorTranslations(nil); err == nil {
		t.Fatal("InitValidatorTranslations expected nil provider error")
	}
}
