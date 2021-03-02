package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/log"

	"charlesbases/http/middleware"
)

func main() {
	defer log.Flush()

	n := negroni.New()

	// middleware
	n.Use(middleware.Cors())
	n.Use(middleware.Recovery())
	n.Use(middleware.Negroni())
	n.Use(middleware.JWT())

	n.UseHandler(router())
	n.Run(":8080")
}

// router router
func router() *mux.Router {
	r := mux.NewRouter()

	// 只匹配 GET | POST
	r.Methods("GET", "POST")

	// websocket
	websocket := r.PathPrefix("/stream").Subrouter()
	websocket.Handle("/", nil)

	// rpc
	rpc := r.PathPrefix("/api").Subrouter()
	rpc.Handle("/{service:[a-zA-Z0-9]+}/{endpoint:[a-zA-Z0-9/]+}", nil)

	return r
}
