package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
)

func registerPingRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("", []router.Route{
		// GET
		router.NewGetRoute("/ping", server.Wrap(h.PingGet)),
		// POST
		router.NewPostRoute("/ping", server.Wrap(h.PingPost)),
	})
}
