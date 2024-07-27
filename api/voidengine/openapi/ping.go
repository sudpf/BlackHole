package openapi

import (
	"BlackHole/api/router"
	"BlackHole/api/wrapper"
	"BlackHole/internal/voidengine/controller"
)

func init() {
	RegisteRoutes("v1", []router.Route{
		// GET
		router.NewGetRoute("/ping", wrapper.WrapperEnvFunc(controller.PingGet)),
		// POST
		router.NewPostRoute("/ping", wrapper.WrapperEnvFunc(controller.PingPost)),
	})
}
