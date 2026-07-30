package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func registerPingRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("", []router.Route{
		// GET
		router.NewGetRoute("/ping", wrapper.WrapperEnvFunc(h.PingGet)),
		// POST
		router.NewPostRoute("/ping", wrapper.WrapperEnvFunc(h.PingPost)),
	})
}
