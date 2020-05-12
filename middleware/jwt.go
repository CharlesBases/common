package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/CharlesBases/common/auth"
	"github.com/CharlesBases/common/log"
)

var JWTConfig *jwtConfig

type (
	jwtConfig struct {
		interceptConfig
		SecretKey   string
		VerifyToken func(token string) bool
	}

	interceptConfig struct {
		Includes []string // 优先级：高
		Excludes []string // 优先级：低
		Fast     bool     // ture:direct，false:regexp
	}
)

func JWT() func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if JWTConfig.intercept(r) {
			token := r.Header.Get("Authorization")

			user, err := auth.ParToken(r)
			if err != nil {
				log.Error(err)
				tokenError(rw)
				return
			}
			if JWTConfig.VerifyToken != nil {
				if !JWTConfig.VerifyToken(token) {
					log.Warn("token logout")
					tokenError(rw)
					return
				}
			}
			next(rw, r.WithContext(context.WithValue(r.Context(), "user_id", user.UserID)))
			return
		}

		next(rw, r)
	}
}

func (config *jwtConfig) intercept(r *http.Request) bool {
	if config == nil {
		return true
	}

	requestURI := strings.Split(r.RequestURI, "?")[0]

	for _, uri := range config.Includes {
		switch config.Fast {
		case true:
			if strings.HasPrefix(requestURI, uri) {
				return true
			}
		case false:
			if regexp.MustCompile(uri).MatchString(requestURI) {
				return true
			}
		}
	}

	for _, uri := range config.Excludes {
		switch config.Fast {
		case true:
			if strings.HasPrefix(requestURI, uri) {
				log.Warn("don't interception:", requestURI)
				return false
			}
		case false:
			if regexp.MustCompile(uri).MatchString(requestURI) {
				log.Warn("don't interception:", requestURI)
				return false
			}
		}
	}

	return true
}

func tokenError(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
	data, _ := json.Marshal(map[string]interface{}{
		"err": 401,
		"msg": "请求错误",
	})
	rw.Write(data)
}
