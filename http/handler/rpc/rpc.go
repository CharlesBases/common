package rpc

import (
	"net/http"
)

// rpc .
type rpc struct{}

// NewHandler .
func NewHandler() *rpc {
	return new(rpc)
}

func (rpc *rpc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API test"))
	return
}
