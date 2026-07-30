package router

import (
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
)

func RegisterRoutes(server *openapi.Server, h *handler.Handler) {
	registerPingRoutes(server, h)
	registerTrafficRoutes(server, h)
	registerUserRoutes(server, h)
}
