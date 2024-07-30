package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func init() {
	openapi.RegisteRoutes("v1", []router.Route{
		router.NewGetRoute("/user", wrapper.WrapperEnvFunc(handler.ListUser)),
		router.NewPostRoute("/user", wrapper.WrapperEnvFunc(handler.AddUer)),
		router.NewPutRoute("/user", wrapper.WrapperEnvFunc(handler.ModifyUer)),
		router.NewDeleteRoute("/user", wrapper.WrapperEnvFunc(handler.DeleteUer)),
	})
}
