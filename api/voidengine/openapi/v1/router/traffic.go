package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/api/wrapper"
)

func init() {
	openapi.RegisteRoutes("v1", []router.Route{
		router.NewGetRoute("/traffic", wrapper.WrapperEnvFunc(handler.ListNetworkTraffic)),
	})
}
