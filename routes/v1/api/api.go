package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/inhzus/go-berater/config"
	"github.com/inhzus/go-berater/middlewares"
	"github.com/inhzus/go-berater/models"
	"github.com/inhzus/go-berater/utils"
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
		auth.GET("/code/:code", checkCode)
		auth.POST("/candidate", addCandidate)
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
	claims := c.MustGet(conf.Jwt.Identity).(*middlewares.OpenidClaims)
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
	genCode := strconv.Itoa(lowerBound + rand.Intn(upperBound-lowerBound))
	err = utils.SendSMS(phoneJson.Phone, genCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error(),})
		return
	}
	err = client.HMSet(claims.Openid+conf.Code.Suffix, map[string]interface{}{
		"code":   genCode,
		"phone":  phoneJson.Phone,
		"status": false,
	}).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error(),})
		return
	}
	client.Expire(claims.Openid+conf.Code.Suffix, time.Duration(conf.Code.ExpireTime)*time.Minute)
	c.Status(http.StatusOK)
}

func checkCode(c *gin.Context) {
	conf := config.GetConfig()
	claims := c.MustGet(conf.Jwt.Identity).(*middlewares.OpenidClaims)
	redisKey := claims.Openid + conf.Code.Suffix
	code := c.Param("code")
	cached, err := client.HGetAll(redisKey).Result()
	if err != nil || cached["code"] != code {
		c.Status(http.StatusNotFound)
		return
	}
	client.HSet(redisKey, "status", true)
	client.Expire(redisKey, time.Duration(conf.Code.ExpireTime)*time.Minute)
	c.Status(http.StatusOK)
}

func addCandidate(c *gin.Context) {
	conf := config.GetConfig()
	openid := c.MustGet(conf.Jwt.Identity).(*middlewares.OpenidClaims).Openid
	redisKey := openid + conf.Code.Suffix
	cached, err := client.HGet(redisKey, "status").Result()
	if status, _ := strconv.ParseBool(cached); err != nil || !status {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Phone not verified",})
	}
	var candidate models.Candidate
	err = c.ShouldBindJSON(&candidate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(),})
		return
	}
	candidate.Openid = openid
	if models.ExistCandidateById(openid) {
		c.JSON(http.StatusConflict, gin.H{"msg": "Candidate has been created with the openid",})
		return
	}
	err = models.AddCandidate(&candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error(),})
	} else {
		c.Status(http.StatusOK)
	}
}
