package openapi

import (
	"BlackHole/api/common/response"
	"BlackHole/api/middleware"
	"BlackHole/api/router"
	"BlackHole/api/swagger"
	"BlackHole/docs/api/voidengine"
	"BlackHole/internal/voidengine/locales"
	"BlackHole/pkg/env"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	address string
	router  *gin.Engine
	routes  map[string][]router.Route
}

func NewHTTPServer(address, apiLogFile string) *Server {
	gin.DefaultWriter = log.StandardLogger().Out
	gin.DefaultErrorWriter = log.StandardLogger().Out

	server := &Server{
		address: address,
		router:  gin.New(),
		routes:  make(map[string][]router.Route),
	}

	env.SetupTranslations()
	env.InitLocalizer(locales.EnTranslations, locales.ZhTranslations)

	middleware.ApiLogMiddlewares(server.router, apiLogFile)

	server.router.NoRoute(func(c *gin.Context) {
		e := env.NewEnv(c.GetHeader("Accept-Language"), c.ClientIP())
		c.JSON(http.StatusNotFound, response.ApiNotFound.Tr(e))
	})

	swagger.SwaggerGenerator(server.router)
	voidengine.SwaggerInfo.Title = "VoidEngen"
	voidengine.SwaggerInfo.Version = "v1"
	voidengine.SwaggerInfo.Description = "API 文档"
	voidengine.SwaggerInfo.Host = swaggerHost(address)
	voidengine.SwaggerInfo.BasePath = "/"
	server.router.Static("/voidengine", "docs/api/voidengine")

	return server
}

func (s *Server) HTTPServer() *http.Server {
	for groupStr, routes := range s.routes {
		group := s.router.Group(groupStr)

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

	return &http.Server{
		Addr:    s.address,
		Handler: s.router,
	}
}

func swaggerHost(address string) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return ""
	}

	switch host {
	case "", "0.0.0.0", "::", "127.0.0.1", "localhost":
		return ""
	}

	return net.JoinHostPort(host, port)
}

func (s *Server) RegisterRoutes(group string, routes []router.Route) {
	s.routes[group] = append(s.routes[group], routes...)
}
