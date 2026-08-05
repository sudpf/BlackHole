package router

import (
	"BlackHole/api/router"
	"BlackHole/api/voidengine/openapi"
	"BlackHole/api/voidengine/openapi/v1/handler"
	"BlackHole/pkg/auth"
)

func registerUserRoutes(server *openapi.Server, h *handler.Handler, authenticator auth.Authenticator) {
	server.RegisterRoutes("v1", []router.Route{
		router.NewGetRoute("/user", server.Wrap(h.ListUser), openapi.RequireAuth(authenticator)),
		router.NewPostRoute("/user", server.Wrap(h.AddUser), openapi.RequireAuth(authenticator)),
		router.NewPutRoute("/user", server.Wrap(h.ModifyUser), openapi.RequireAuth(authenticator)),
		router.NewDeleteRoute("/user", server.Wrap(h.DeleteUser), openapi.RequireAuth(authenticator)),
	})
}
