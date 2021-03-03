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

type (
	jwtConfig struct {
		interceptConfig
		SecretKey   string
		VerifyToken func(token string) bool
	}

	interceptConfig struct {
		Includes []string // 优先级：高
		Excludes []string // 优先级：低
		Prefix   string   // 前缀
		Fast     bool     // ture:direct，false:regexp
	}
)

func JWT() *jwtConfig {
	return new(jwtConfig)
}

func (config *jwtConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if config.intercept(r) {
		token := r.Header.Get("Authorization")

		user, err := auth.ParTokenFromRequest(r)
		if err != nil {
			log.Error(err)
			tokenError(rw)
			return
		}
		if config.VerifyToken != nil {
			if !config.VerifyToken(strings.TrimPrefix(token, config.Prefix)) {
				log.Warn("token invalid")
				tokenError(rw)
				return
			}
		}
		next(rw, r.WithContext(context.WithValue(r.Context(), "user_id", user.UserID)))
		return
	}

	next(rw, r)
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
		default:
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
		default:
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
		"code":    http.StatusUnauthorized,
		"message": http.StatusText(http.StatusUnauthorized),
	})
	rw.Write(data)
}
