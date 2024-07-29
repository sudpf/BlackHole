package v1

import (
	"BlackHole/api/openapi"
	"BlackHole/api/router"
	"BlackHole/api/wrapper"
	"BlackHole/internal/controller/voidengine"
)

func init() {
	openapi.RegisteRoutes("v1", []router.Route{
		// GET
		router.NewGetRoute("/ping", wrapper.WrapperEnvFunc(voidengine.PingGet)),
		// POST
		router.NewPostRoute("/ping", wrapper.WrapperEnvFunc(voidengine.PingPost)),
	})
}
