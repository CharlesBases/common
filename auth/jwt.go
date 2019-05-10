package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/log"
)

type TokenPayload struct {
	UserId    int   `json:"userId"`
	Timestamp int64 `json:"timestamp"`
}

// jwt config
type JWTConfig struct {
	InterceptConfig                                               // 过滤规则
	SecretKey         string                                      // 密钥
	CheckTokenPayload func(token string, load *TokenPayload) bool // 验证Token
}

type InterceptConfig interface {
	Includes() []string
	Excludes() []string
	Fast() bool // ture:direct，false:regexp
}

type DefaultInterceptConfig struct {
	Include []string
	Exclude []string
	IsFast  bool
}

func (de DefaultInterceptConfig) Includes() []string {
	return de.Include
}

func (de DefaultInterceptConfig) Excludes() []string {
	return de.Exclude
}

func (de DefaultInterceptConfig) Fast() bool {
	return de.IsFast
}

type FastInterceptConfig struct {
	Include []string
	Exclude []string
}

func (fa FastInterceptConfig) Includes() []string {
	return fa.Include
}

func (fa FastInterceptConfig) Excludes() []string {
	return fa.Exclude
}

func (fa FastInterceptConfig) Fast() bool {
	return true
}

func JWT(jwtcfg JWTConfig) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if strings.ToUpper(r.Method) == "OPTIONS" {
			next(w, r)
			return
		}
		// log.Debug("Authorization:", r.Header["Authorization"])
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtcfg.SecretKey), nil
			})
		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			var user = claims["user"]
			if user != nil {
				var load = new(TokenPayload)
				bytes, _ := json.Marshal(user)
				json.Unmarshal(bytes, load)
				if jwtcfg.CheckTokenPayload != nil {
					if !jwtcfg.CheckTokenPayload(token.Raw, load) {
						log.Warn("token logout")
						w.WriteHeader(http.StatusUnauthorized)
						result := map[string]interface{}{
							"code": 999,
							"flag": false,
							"msg":  "Token无效",
						}
						b, _ := json.Marshal(result)
						w.Write(b)
						return
					}
				}
				ctx := r.Context()
				ctx = context.WithValue(ctx, "userId", load.UserId)
				ctx = context.WithValue(ctx, "user", *load)
				request := r.WithContext(ctx)
				next(w, request)
				return
			}
		}
		switch err.(type) {
		case *jwt.ValidationError:
			log.Warn(r.RequestURI, " token.Valid: ", err)
		default:
			log.Error(r.RequestURI, " token.Valid: ", err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		result := map[string]interface{}{
			"code": 999,
			"flag": false,
			"msg":  "Token无效",
		}
		b, _ := json.Marshal(result)
		w.Write(b)
	}
	return InterceptHandlerFunc(jwtcfg, h)
}

// 拦截配置
func InterceptHandlerFunc(cfg InterceptConfig, h negroni.HandlerFunc) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if cfg == nil {
			h(rw, r, next)
			return
		}
		excludes := cfg.Excludes()
		includes := cfg.Includes()
		var pass bool
		requestURI := r.RequestURI
		index := strings.Index(requestURI, "?")
		if index != -1 {
			requestURI = requestURI[0:index]
		}
		if includes == nil || len(includes) == 0 {
			pass = true
		} else {
			for i := range includes {
				if cfg.Fast() {
					if includes[i] == requestURI {
						pass = true
						break
					}
				} else {
					matched, err := regexp.MatchString(includes[i], requestURI)
					if err != nil {
						log.Error(err)
						continue
					}
					if matched {
						pass = matched
						break
					}
				}
			}
		}
		if !pass {
			next(rw, r)
			return
		}
		if excludes == nil || len(excludes) == 0 {
			pass = true
		} else {
			for i := range excludes {
				if cfg.Fast() {
					if excludes[i] == requestURI {
						pass = false
						break
					}
				} else {
					matched, err := regexp.MatchString(excludes[i], requestURI)
					if err != nil {
						log.Error(err)
						continue
					}
					if matched {
						pass = !matched
						break
					}
				}
			}
		}
		if pass {
			h(rw, r, next)
		} else {
			log.Debug("不拦截：", requestURI)
			next(rw, r)
		}
	}
}
