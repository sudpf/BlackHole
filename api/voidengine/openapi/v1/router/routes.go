package router

import "BlackHole/api/voidengine/openapi/v1/handler"

func RegisterRoutes(h *handler.Handler) {
	registerPingRoutes(h)
	registerTrafficRoutes(h)
	registerUserRoutes(h)
}
