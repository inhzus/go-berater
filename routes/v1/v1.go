package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/inhzus/go-berater/routes/v1/api"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ApplyRoutes(group *gin.RouterGroup) {
	v1 := group.Group("/v1")
	{
		v1.GET("/ping", ping)
		api.ApplyRoutes(v1)
	}
}
