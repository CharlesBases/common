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

var (
	tokenError = map[string]interface{}{
		"code": 999,
		"flag": false,
		"msg":  "token is not valid",
	}
)

type User struct {
	UserId    int   `json:"userId"`
	Timestamp int64 `json:"timestamp"`
}

type jwtClaims struct {
	User User
	jwt.StandardClaims
}

// jwt config
type JWTConfig struct {
	InterceptConfig                                     // 过滤规则
	SecretKey       string                              // 密钥
	VerifyToken     func(token string, user *User) bool // 验证Token
}

type InterceptConfig interface {
	Includes() []string // 优先级：高
	Excludes() []string // 优先级：低
	Fast() bool         // ture:direct，false:regexp
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
		if err == nil {
			if token.Valid {
				user := new(User)
				bytes, _ := json.Marshal(token.Claims.(jwt.MapClaims)["user"])
				json.Unmarshal(bytes, user)
				if jwtcfg.VerifyToken != nil {
					if !jwtcfg.VerifyToken(token.Raw, user) {
						log.Warn("token logout")
						w.WriteHeader(http.StatusUnauthorized)
						b, _ := json.Marshal(tokenError)
						w.Write(b)
						return
					}
				}
				ctx := context.WithValue(r.Context(), "userId", user.UserId)
				request := r.WithContext(ctx)
				next(w, request)
				return
			} else {
				log.Warn(r.RequestURI, " - token is not valid")
			}
		}
		log.Error(r.RequestURI, " - unauthorized access to this resource:  ", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(&tokenError)
		w.Write(b)
	}
	return intercept(jwtcfg, h)
}

// 拦截配置
func intercept(cfg InterceptConfig, h negroni.HandlerFunc) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if cfg == nil {
			h(rw, r, next)
			return
		}
		isVerify := true
		requestURI := r.RequestURI
		if index := strings.Index(requestURI, "?"); index != -1 {
			requestURI = requestURI[:index]
		}

		if len(cfg.Includes()) != 0 {
			isVerify = false
			for _, v := range cfg.Includes() {
				if cfg.Fast() {
					if v == requestURI {
						isVerify = true
						break
					}
				} else {
					if regexp.MustCompile(v).MatchString(requestURI) {
						isVerify = true
						break
					}
				}
			}
		}
		if len(cfg.Excludes()) != 0 {
			for _, v := range cfg.Excludes() {
				if cfg.Fast() {
					if v == requestURI {
						isVerify = false
						break

					}
				} else {
					if regexp.MustCompile(v).MatchString(requestURI) {
						isVerify = false
						break
					}
				}
			}
		}
		if isVerify {
			h(rw, r, next)
		} else {
			log.Warn("don't interception: ", requestURI)
			next(rw, r)
		}
	}
}
