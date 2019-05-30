package api

import (
	"github.com/gin-gonic/gin"
	"github.com/inhzus/go-berater/middlewares/jwt"
	"net/http"
)

func ApplyRoutes(r *gin.RouterGroup) {
	api := r.Group("/api")
	{
		api.GET("/test/token/:openid", testToken)
	}
	auth := api.Group("")
	auth.Use(jwt.JwtMiddleware())
	{
		auth.GET("/token", checkToken)
	}
}

func testToken(c *gin.Context) {
	openid := c.Param("openid")
	auth, err := jwt.CreateToken(openid)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": auth,
	})
}

func checkToken(c *gin.Context) {
	c.Status(http.StatusOK)
}
