package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (pr *pingRouter) pingGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get pong",
	})
}

func (pr *pingRouter) pingPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "post pong",
	})
}
