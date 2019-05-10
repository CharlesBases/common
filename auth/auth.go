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

func (payload TokenPayload) GenRedisKey(prefix string) string {
	return fmt.Sprintf("%s%s_%d", prefix, strconv.Itoa(payload.UserId), payload.Timestamp)
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

func GenToken(secretKey string, duration time.Duration, payload *TokenPayload) (tokenStr string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user"] = map[string]interface{}{"user_id": payload.UserId, "timestamp": payload.Timestamp}
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func GetToken(conn redis.Conn, redisKey string) (tokenStr string, err error) {
	return redis.String(conn.Do("GET", redisKey))
}

func memoryToken(conn redis.Conn, redisKey string, value string) error {
	_, err := conn.Do("SET", redisKey, value)
	return err
}
