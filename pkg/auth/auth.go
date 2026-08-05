package auth

import (
	"context"
	"net/http"
)

type Method string

const (
	MethodIAMJWT   Method = "iam_jwt"
	MethodAPIKey   Method = "api_key"
	MethodPassword Method = "password"
	MethodLocalJWT Method = "local_jwt"
)

type Principal struct {
	Method  Method
	Subject string
	UserID  string
	Name    string
	Tenant  string
	Scopes  []string
	Roles   []string
	Claims  map[string]any
}

type Request struct {
	Authorization string
	APIKey        string
	Header        http.Header
	Method        string
	Path          string
}

type Authenticator interface {
	Authenticate(ctx context.Context, request Request) (*Principal, error)
}

func RequestFromHTTP(r *http.Request) Request {
	if r == nil {
		return Request{}
	}

	return Request{
		Authorization: r.Header.Get("Authorization"),
		APIKey:        r.Header.Get("X-API-Key"),
		Header:        r.Header.Clone(),
		Method:        r.Method,
		Path:          r.URL.Path,
	}
}
