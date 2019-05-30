package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/inhzus/go-berater/config"
	"net/http"
	"strings"
	"time"
)

var (
	TokenInvalid = errors.New("token invalid")
	SecretKey    = []byte( config.GetConfig().Jwt.SecretKey)
)

type OpenidClaims struct {
	Openid string
	jwt.StandardClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		claims, err := ParseToken(auth)
		if err != nil {
			_ = context.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		context.Set("claims", claims)
	}
}

func CreateToken(openid string) (string, error) {
	c := config.GetConfig()
	claims := OpenidClaims{
		openid,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Local().Add(time.Duration(c.Jwt.Timeout) * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ParseToken(auth string) (*OpenidClaims, error) {
	authArray := strings.Split(auth, " ")
	if len(authArray) < 2 {
		return nil, TokenInvalid
	}
	auth = authArray[1]
	token, err := jwt.ParseWithClaims(auth, &OpenidClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, TokenInvalid
	}
	if claims, status := token.Claims.(*OpenidClaims); status && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
