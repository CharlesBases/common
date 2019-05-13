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

func GetUser(r *http.Request) (userId int, err error) {
	value := r.Context().Value("userId")
	userId, err = strconv.Atoi(fmt.Sprintf(`%v`, value))
	if err != nil {
		return -1, errors.New("userId has no value in http request context")
	}
	return userId, nil
}

func GenToken(SecretKey string, duration time.Duration, load *User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
		"user": load,
	})).SignedString([]byte(SecretKey))
}

func ParToken(SecretKey string, tokenString string) (*User, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return &claims.User, nil
	}
	return nil, errors.New("token is error")
}

func (user *User) GenRedisKey(prefix string) string {
	return fmt.Sprintf("%s%s_%d", prefix, strconv.Itoa(user.UserId), user.Timestamp)
}

func SetToken(conn redis.Conn, redisKey string, value string) error {
	_, err := conn.Do("SET", redisKey, value)
	return err
}

func GetToken(conn redis.Conn, redisKey string) (tokenStr string, err error) {
	return redis.String(conn.Do("GET", redisKey))
}
