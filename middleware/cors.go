package middleware

import (
	"net/http"

	"charlesbases/common/log"
)

const (
	allowOrigin      = "*"
	allowMethods     = "GET, POST, OPTIONS, PUT, PATCH, DELETE"
	allowHeaders     = "Authorization, Content-Length, X-CSRF-Token, Token,session, X_Requested_With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma"
	allowCredentials = "false"
	exposeHeaders    = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma, FooBar"
	maxAge           = "86400"
	contentType      = "application/json"
)

func Cors() func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.Method == "OPTIONS" {
			log.Info(r.RequestURI, " Method:OPTIONS")
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Header.Get("Origin") != "" {
			// 允许访问所有域
			rw.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			// 支持跨域的请求方法
			rw.Header().Set("Access-Control-Allow-Methods", allowMethods)
			// 支持跨域的header的类型
			rw.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			// cookie
			rw.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
			// 跨域关键设置，让浏览器可以解析
			rw.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			// 缓存
			rw.Header().Set("Access-Control-Max-Age", maxAge)
			// 返回格式
			rw.Header().Set("content-type", contentType)
		}
		next(rw, r)
	}
}
