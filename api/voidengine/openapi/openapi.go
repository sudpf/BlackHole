package openapi

import (
	"BlackHole/api/middleware"
	"BlackHole/api/router"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

var (
	apiRouter *gin.Engine
	apiRoutes = make(map[string][]router.Route)
)

func Run() {
	apiRouter.Run(":8080")
}

func InitApi() {
	apiRouter = gin.New()

	middleware.InitMiddlewares(apiRouter)

	for groupStr, routes := range apiRoutes {
		group := apiRouter.Group(groupStr)

		for _, route := range routes {
			switch route.Method() {
			case http.MethodGet:
				group.GET(route.Path(), route.Handler())
			case http.MethodHead:
				group.HEAD(route.Path(), route.Handler())
			case http.MethodPost:
				group.POST(route.Path(), route.Handler())
			case http.MethodPut:
				group.PUT(route.Path(), route.Handler())
			case http.MethodPatch:
				group.PATCH(route.Path(), route.Handler())
			case http.MethodDelete:
				group.DELETE(route.Path(), route.Handler())
			case http.MethodOptions:
				group.OPTIONS(route.Path(), route.Handler())
			default:
				log.WithField("method", route.Method()).Error("unknown method")
			}
		}
	}
}

func RegisteRoutes(group string, routes []router.Route) {
	apiRoutes[group] = append(apiRoutes[group], routes...)
}
