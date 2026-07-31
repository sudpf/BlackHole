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
		Language: constant.LangChinese,
		ClientIP: "127.0.0.1",
	})

	ev := provider.NewEnvFromContext(ctx)
	if ev.RequestId != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("RequestId = %q, want 4bf92f3577b34da6a3ce929d0e0e4736", ev.RequestId)
	}
	if ev.ClientIp != "127.0.0.1" {
		t.Fatalf("ClientIp = %q, want 127.0.0.1", ev.ClientIp)
	}
	if got := ev.MustLocalize("hello"); got != "你好" {
		t.Fatalf("MustLocalize = %q, want 你好", got)
	}
}
