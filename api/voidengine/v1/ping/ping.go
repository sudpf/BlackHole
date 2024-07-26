package ping

import "BlackHole/api/router"

type pingRouter struct {
	routes []router.Route
}

// NewRouter initializes a new image router
func NewRouter() router.Router {
	pr := &pingRouter{}
	pr.initRoutes()
	return pr
}

// Routes returns the available routes to the image controller
func (pr *pingRouter) Routes() []router.Route {
	return pr.routes
}

// initRoutes initializes the routes in the image router
func (pr *pingRouter) initRoutes() {
	pr.routes = []router.Route{
		// GET
		router.NewGetRoute("/ping", pr.pingGet),
		// POST
		router.NewPostRoute("/ping", pr.pingPost),
	}
}
