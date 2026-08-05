package router

import (
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/pkg/auth"
)

type Options struct {
	Authenticator auth.Authenticator
}

func RegisterRoutes(server *openapi.Server, h *handler.Handler, opts ...Options) {
	options := Options{}
	if len(opts) > 0 {
		options = opts[0]
	}

	registerPingRoutes(server, h)
	registerTrafficRoutes(server, h)
	registerUserRoutes(server, h, options.Authenticator)
}
