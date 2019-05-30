package routes

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/inhzus/go-berater/routes/v1"
)

func ApplyRoutes(r *gin.Engine) {
	router := r.Group("/")
	{
		v1.ApplyRoutes(router)
	}
}
