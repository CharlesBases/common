package middleware

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/CharlesBases/common/log"
	"github.com/CharlesBases/common/web/weberror"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("token's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	SignKey          = "newtrekWang"
)

type jwtSign struct {
	SigningKey []byte
}

type user struct {
	ID int `json:"userId"`
	jwt.StandardClaims
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			log.Error(errors.New("无权限访问"))
			c.JSON(http.StatusOK, weberror.NewWebrror(weberror.TokenNotFound))
			c.Abort()
			return
		}
	}
}

func (j *jwtSign) parseToken(tokenString string) (*user, error) {
	token, err := jwt.ParseWithClaims(tokenString, &user{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*user); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
