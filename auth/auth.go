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
	if value != nil {
		f, ok := value.(float64)
		if ok {
			userId = int(f)
		} else {
			userId = value.(int)
		}
		return userId, nil
	} else {
		return userId, errors.New("userId has no value in http request context")
	}
}

func GenToken(secretKey string, duration time.Duration, load *TokenPayload) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
		"user": load,
	})).SignedString([]byte(securityKey))
}

func SetToken(conn redis.Conn, redisKey string, value string) error {
	_, err := conn.Do("SET", redisKey, value)
	return err
}

func GetToken(conn redis.Conn, redisKey string) (tokenStr string, err error) {
	return redis.String(conn.Do("GET", redisKey))
}
