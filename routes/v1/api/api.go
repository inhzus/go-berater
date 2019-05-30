package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/inhzus/go-berater/config"
	"github.com/inhzus/go-berater/middlewares"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	client = redis.NewClient(&redis.Options{
		Addr: config.GetConfig().Redis.Addr,
	})
)

type CodeStorage struct {
	Code   string
	phone  string
	Status bool
}

func ApplyRoutes(r *gin.RouterGroup) {
	api := r.Group("/api")
	{
		api.GET("/test/token/:openid", testToken)
	}
	auth := api.Group("")
	auth.Use(middlewares.JwtMiddleware())
	{
		auth.GET("/token", checkToken)
		auth.POST("/code", sendCode)
	}
}

func testToken(c *gin.Context) {
	openid := c.Param("openid")
	auth, err := middlewares.CreateToken(openid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(),})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": auth,
	})
}

func checkToken(c *gin.Context) {
	c.Status(http.StatusOK)
}

func sendCode(c *gin.Context) {
	conf := config.GetConfig()
	claims := c.MustGet("claims").(*middlewares.OpenidClaims)
	var phoneJson struct {
		Phone string
	}
	err := c.ShouldBindJSON(&phoneJson)
	//err := c.BindJSON(&phoneJson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Request JSON format error or \"phone\" missing",
		})
		return
	}
	lowerBound := 1
	upperBound := 10
	for i := 1; i != conf.Code.Length; i++ {
		lowerBound *= 10
		upperBound *= 10
	}
	genCode := strconv.Itoa(lowerBound + rand.Intn(upperBound - lowerBound))
	err = client.HMSet(claims.Openid+conf.Code.Suffix, map[string]interface{}{
		"code":   genCode,
		"phone":  phoneJson.Phone,
		"status": false,
	}).Err()
	client.Expire(claims.Openid+conf.Code.Suffix, time.Duration(conf.Code.ExpireTime)*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error(),})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": genCode,
	})
}
