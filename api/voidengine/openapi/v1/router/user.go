package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func registerUserRoutes() {
	openapi.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/user", wrapper.WrapperEnvFunc(handler.ListUser)),
		router.NewPostRoute("/user", wrapper.WrapperEnvFunc(handler.AddUser)),
		router.NewPutRoute("/user", wrapper.WrapperEnvFunc(handler.ModifyUser)),
		router.NewDeleteRoute("/user", wrapper.WrapperEnvFunc(handler.DeleteUser)),
	})
}
