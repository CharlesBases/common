package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CharlesBases/common/log"
)

func Cors() func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		origin := r.Header.Get("Origin")

		var headerKeys []string
		for k, _ := range r.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			// 允许访问所有域
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			// 支持跨域的请求方法
			rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// 支持跨域的header的类型
			rw.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// cookie
			rw.Header().Set("Access-Control-Allow-Credentials", "false")
			// 跨域关键设置，让浏览器可以解析
			rw.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			// 缓存
			rw.Header().Set("Access-Control-Max-Age", "86400")
			// 返回格式
			rw.Header().Set("content-type", "application/json")
		}
		if r.Method == "OPTIONS" {
			log.Info(r.RequestURI, "OPTIONS request")
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		next(rw, r)
	}
}
