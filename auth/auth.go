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
	value := r.FormValue("userId")
	userId, err = strconv.Atoi(fmt.Sprintf(`%v`, value))
	if err != nil {
		return -1, errors.New("userId has no value in http request context")
	}
	return userId, nil
}

// generate token
func GenToken(SecretKey string, duration time.Duration, user *User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(map[string]interface{}{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
		"user": user,
	})).SignedString([]byte(SecretKey))
}

// parse token
func ParToken(SecretKey string, tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err == nil {
		if token.Valid {
			return &token.Claims.(*jwtClaims).User, nil
		} else {
			return nil, errors.New("token is not valid")
		}
	}
	return nil, errors.New("unauthorized access to this resource")
}

/*
 storage token in redis
*/
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
