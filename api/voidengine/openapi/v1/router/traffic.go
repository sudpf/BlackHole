package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
)

func registerTrafficRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/traffic", server.Wrap(h.ListNetworkTraffic)),
	})
}
