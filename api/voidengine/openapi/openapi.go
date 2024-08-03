package openapi

import (
	"BlackHole/api/middleware"
	"BlackHole/api/router"
	"BlackHole/api/swagger"
	"BlackHole/docs/api/voidengine"
	"BlackHole/internal/voidengine/locales"
	"BlackHole/internal/voidengine/response"
	"BlackHole/pkg/env"
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
	gin.DefaultWriter = log.StandardLogger().Out
	gin.DefaultErrorWriter = log.StandardLogger().Out

	apiRouter = gin.New()

	env.SetupTranslations()
	env.InitLocalizer(locales.EnTranslations, locales.ZhTranslations)

	middleware.ApiLogMiddlewares(apiRouter)

	apiRouter.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.InvalidParams)
	})

	swagger.SwaggerGenerator(apiRouter)
	voidengine.SwaggerInfo.Title = "VoidEngen"
	voidengine.SwaggerInfo.Version = "v1"
	voidengine.SwaggerInfo.Description = "API 文档"
	voidengine.SwaggerInfo.Host = "127.0.0.1:8080"
	voidengine.SwaggerInfo.BasePath = "/"
	apiRouter.Static("/voidengine", "docs/api/voidengine")

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
