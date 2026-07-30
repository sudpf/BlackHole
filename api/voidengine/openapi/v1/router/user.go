package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func registerUserRoutes(server *openapi.Server, h *handler.Handler) {
	server.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/user", wrapper.WrapperEnvFunc(h.ListUser)),
		router.NewPostRoute("/user", wrapper.WrapperEnvFunc(h.AddUser)),
		router.NewPutRoute("/user", wrapper.WrapperEnvFunc(h.ModifyUser)),
		router.NewDeleteRoute("/user", wrapper.WrapperEnvFunc(h.DeleteUser)),
	})
}
