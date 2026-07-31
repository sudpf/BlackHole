package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
)

func registerUserRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/user", server.WrapEnv(h.ListUser)),
		router.NewPostRoute("/user", server.WrapEnv(h.AddUser)),
		router.NewPutRoute("/user", server.WrapEnv(h.ModifyUser)),
		router.NewDeleteRoute("/user", server.WrapEnv(h.DeleteUser)),
	})
}
