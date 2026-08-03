package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
)

func registerUserRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/user", server.Wrap(h.ListUser)),
		router.NewPostRoute("/user", server.Wrap(h.AddUser)),
		router.NewPutRoute("/user", server.Wrap(h.ModifyUser)),
		router.NewDeleteRoute("/user", server.Wrap(h.DeleteUser)),
	})
}
