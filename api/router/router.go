package router

import "github.com/gin-gonic/gin"

type Router interface {
	Routes() []Route
}

type Route interface {
	Handler() gin.HandlerFunc
	Method() string
	Path() string
}
