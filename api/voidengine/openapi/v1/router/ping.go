package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func init() {
	openapi.RegisteRoutes("", []router.Route{
		// GET
		router.NewGetRoute("/ping", wrapper.WrapperEnvFunc(handler.PingGet)),
		// POST
		router.NewPostRoute("/ping", wrapper.WrapperEnvFunc(handler.PingPost)),
	})
}
