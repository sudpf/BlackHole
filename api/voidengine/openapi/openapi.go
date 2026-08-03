package openapi

import (
	"BlackHole/api/middleware"
	"BlackHole/api/router"
	"BlackHole/api/swagger"
	"BlackHole/api/wrapper"
	"BlackHole/docs/api/voidengine"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/locales"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/env"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	address      string
	router       *gin.Engine
	routes       map[string][]router.Route
	envProvider  *env.Provider
	errorCatalog *apperror.Catalog
}

func NewHTTPServer(address, apiLogFile string, apiLogSize string, requestTimeout time.Duration) (*Server, error) {
	envProvider, err := env.NewProvider(locales.English(), locales.Chinese())
	if err != nil {
		return nil, fmt.Errorf("initialize env provider: %w", err)
	}
	if err = env.InitValidatorTranslations(envProvider); err != nil {
		return nil, fmt.Errorf("initialize validator translations: %w", err)
	}

	errorCatalog, err := newErrorCatalog()
	if err != nil {
		return nil, fmt.Errorf("initialize error catalog: %w", err)
	}
	if err := envProvider.ValidateMessages(errorCatalog.MessageIDs()); err != nil {
		return nil, fmt.Errorf("validate error messages: %w", err)
	}

	server := &Server{
		address:      address,
		router:       gin.New(),
		routes:       make(map[string][]router.Route),
		envProvider:  envProvider,
		errorCatalog: errorCatalog,
	}

	if err := middleware.ApiLogMiddlewares(server.router, apiLogFile, apiLogSize); err != nil {
		return nil, err
	}
	server.router.Use(middleware.ErrorHandler(envProvider, errorCatalog, errorcode.SystemError))
	server.router.Use(middleware.Recovery(errorcode.SystemError))
	server.router.Use(middleware.RequestContext(requestTimeout))

	server.router.NoRoute(func(c *gin.Context) {
		_ = c.Error(apperror.New(errorcode.APINotFound))
		c.Abort()
	})

	swagger.SwaggerGenerator(server.router)
	voidengine.SwaggerInfo.Title = "VoidEngen"
	voidengine.SwaggerInfo.Version = "v1"
	voidengine.SwaggerInfo.Description = "API 文档"
	voidengine.SwaggerInfo.Host = swaggerHost(address)
	voidengine.SwaggerInfo.BasePath = "/"
	server.router.Static("/voidengine", "docs/api/voidengine")

	return server, nil
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

func (s *Server) Wrap(handler wrapper.HandlerFunc) gin.HandlerFunc {
	return wrapper.Adapt(s.envProvider, handler)
}
