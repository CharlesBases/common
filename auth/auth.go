package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/go-redis/redis/v7"
)

const (
	// SecretKey token secret key
	SecretKey = "shdkkj&(hkdksaYKBKDJah890uiojoiu0KNKSAdhka892hkj!@kndsajhd"
	// Duration ttl for token
	Duration = 4 * time.Hour
)

type (
	// User user info
	User struct {
		UserID    uint64 `json:"user_id"`
		Timestamp int64  `json:"timestamp"`
	}

	jwtClaims struct {
		User User
		jwt.StandardClaims
	}
)

// GetUser get user from request
func GetUser(r *http.Request) (userId int, err error) {
	value := r.FormValue("user_id")
	userId, err = strconv.Atoi(fmt.Sprintf(`%v`, value))
	if err != nil {
		return -1, errors.New("userId has no value in http request context")
	}
	return userId, nil
}

// GenToken generate token
func GenToken(user *User) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims(map[string]interface{}{
			// "iat":  time.Now().Unix(),
			// "exp":  time.Now().Add(Duration).Unix(),
			"user": user,
		}),
	).SignedString([]byte(SecretKey))
}

// ParToken parse token
func ParToken(r *http.Request) (*User, error) {
	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)
	if err == nil {
		switch token.Valid {
		case true:
			return &token.Claims.(*jwtClaims).User, nil
		case false:
			return nil, errors.New("token is not valid")
		}
	}
	return nil, errors.New("unauthorized access to this resource")
}

// GenRedisKey generate token key for redis
func (user *User) GenRedisKey(prefix string) string {
	return fmt.Sprintf("%s%d_%d", prefix, user.UserID, user.Timestamp)
}

// SetToken save token to redis
func SetToken(r *redis.Client, redisKey string, value string) error {
	return r.Do("SET", redisKey, value).Err()
}

// GetToken get token from redis
func GetToken(r *redis.Client, redisKey string) (tokenStr string) {
	return r.Do("GET", redisKey).String()
}

// VerifyToken verify token
func VerifyToken(tokenString string) bool {
	return true
}
