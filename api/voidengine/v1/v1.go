package v1

import (
	"BlackHole/api/middleware"
	"BlackHole/api/voidengine/v1/ping"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type ApiRouterInterface interface {
	GetName() string
	GetUri() string
	GetRouter() gin.IRouter
	SetRouter()
}

type ApiRouter struct {
	Name   string
	Uri    string
	Router gin.IRouter
}

var (
	router     *gin.Engine
	apiRouters map[string]ApiRouterInterface
)

func (a *ApiRouter) GetName() string {
	return a.Name
}
func (a *ApiRouter) GetUri() string {
	return a.Uri
}
func (a *ApiRouter) GetRouter() gin.IRouter {
	return a.Router
}

func (a *ApiRouter) SetRouter() {
	a.Router = router.Group(a.GetUri())
}

func Run() {
	router.Run(":8080")
}

func init() {
	router = gin.New()
	apiRouters = map[string]ApiRouterInterface{}

	routes := ping.NewRouter().Routes()
	for _, route := range routes {
		switch route.Method() {
		case http.MethodGet:
			router.GET(route.Path(), route.Handler())
		}
	}
}

func InitApi() {
	middleware.InitMiddlewares(router)
}

func RegisteRouter(a ApiRouterInterface) {
	if _, ok := apiRouters[a.GetName()]; ok {
		return
	}

	a.SetRouter()
	//a.InitRoute()
	log.WithField("name", a.GetName()).Debug("register api")
	apiRouters[a.GetName()] = a
}
