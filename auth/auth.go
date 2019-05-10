package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
)

func (load TokenPayload) GenRedisKey(prefix string) string {
	return fmt.Sprintf("%s%s_%d", prefix, strconv.Itoa(load.UserId), load.Timestamp)
}

func GetUser(r *http.Request) (userId int, err error) {
	value := r.Context().Value("userId")
	userId, err = strconv.Atoi(fmt.Sprintf(`%v`, value))
	if err != nil {
		return -1, errors.New("userId has no value in http request context")
	}
	return userId, nil
}

// func GetUser(r *http.Request) (userId int, err error) {
// token, _ := jwt.ParseWithClaims(tokenString, &infor{}, func(token *jwt.Token) (interface{}, error) {
// 	return []byte(securityKey), nil
// })
// if claims, ok := token.Claims.(*infor); ok && token.Valid {
// 	return &claims.Infor, nil
// }
// return nil, errors.New("token is error")
// }

func GenToken(sign string, duration time.Duration, load *TokenPayload) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
		"user": load,
	})).SignedString([]byte(sign))
}

func SetToken(conn redis.Conn, redisKey string, value string) error {
	_, err := conn.Do("SET", redisKey, value)
	return err
}

func GetToken(conn redis.Conn, redisKey string) (tokenStr string, err error) {
	return redis.String(conn.Do("GET", redisKey))
}
