package env

import (
	"BlackHole/pkg/auth"
	"BlackHole/pkg/constant"
	"BlackHole/pkg/requestctx"
	"context"
)

type Env struct {
	Lang      string
	ClientIp  string
	RequestId string
	Principal *auth.Principal
}

type contextKey struct{}

func New(lang string, clientIp string) *Env {
	return &Env{Lang: lang, ClientIp: clientIp}
}

func NewFromContext(ctx context.Context) *Env {
	lang := requestctx.Language(ctx)
	if lang == "" {
		lang = constant.LangEnglish
	}

	env := New(lang, requestctx.ClientIP(ctx))
	env.RequestId = requestctx.TraceID(ctx)
	return env
}

func WithContext(ctx context.Context, ev *Env) context.Context {
	return context.WithValue(ctx, contextKey{}, ev)
}

func FromContext(ctx context.Context) (*Env, bool) {
	ev, ok := ctx.Value(contextKey{}).(*Env)
	return ev, ok && ev != nil
}

func (ev *Env) WithPrincipal(principal *auth.Principal) *Env {
	if ev == nil {
		return nil
	}
	copied := *ev
	copied.Principal = principal
	return &copied
}
