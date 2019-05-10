package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/auth"
	"github.com/CharlesBases/common/db"
	"github.com/CharlesBases/common/log"
)

const (
	RedisTokenPrefix = "token_"
	SecurityKey      = "Mdaf43#$%+*/=07RbGc7xkh3frwdIUYknskIHNnJc6_0K240654CCNMmDzxJwVzlPtmcaeadinotii28jWPE"
)

var jwtConfig = auth.JWTConfig{
	InterceptConfig: auth.FastInterceptConfig{
		Exclude: []string{
			"/token",
		}},
	SecretKey: SecurityKey,
	CheckTokenPayload: func(token string, payload *auth.TokenPayload) bool {
		redisKey := payload.GenRedisKey(RedisTokenPrefix)
		tokenStr, err := auth.GetToken(db.InitRedis("127.0.0.1:4399"), redisKey)
		if err != nil {
			log.Warn(err)
		}
		return token == tokenStr
	},
}

func Test(t *testing.T) {
	// router := mux.NewRouter().PathPrefix("/api/").Subrouter()

	defer log.Flush()

	router := gin.New()
	// router
	router.GET("/token", func(c *gin.Context) {
		token, err := auth.GenToken(SecurityKey, time.Hour*4, &auth.TokenPayload{UserId: 1, Timestamp: time.Now().UnixNano()})
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"tiken": token,
		})
	})
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "hahaha",
		})
	})

	n := negroni.New()
	n.Use(Recovery())                                   // recovery
	n.Use(NegroniLogger())                              // logger
	n.Use(NegroniOpentracer())                          // opentracer
	n.UseFunc(negroni.HandlerFunc(Cors()))              // cors
	n.UseFunc(negroni.HandlerFunc(auth.JWT(jwtConfig))) // jwt

	n.UseHandler(router)
	n.Run(":8000")
}
