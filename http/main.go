package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/auth"
	"github.com/CharlesBases/common/log"

	"charlesbases/http/handler/rpc"
	"charlesbases/http/handler/websocket"
	"charlesbases/http/middleware"
)

// defaultPrivateKey default private key for the auth
const defaultPrivateKey = "5aOu5ZOJ5oiR5aSn5Lit5Y2O"

func main() {
	defer log.Flush()

	// init auth
	auth.InitAuth(auth.WithPrivateKey(defaultPrivateKey))

	n := negroni.New()

	// middleware
	n.Use(middleware.Cors())
	n.Use(middleware.Recovery())
	n.Use(middleware.Negroni())
	// n.Use(middleware.JWT())

	n.UseHandler(router())
	n.Run(":8080")
}

// router router 。
func router() *mux.Router {
	r := mux.NewRouter()

	// websocket
	r.Handle("/stream", websocket.NewHandler()).Methods("GET")

	// rpc
	rpcRouter := r.PathPrefix("/api").Subrouter()
	rpcRouter.Handle("{service:[a-zA-Z0-9]+}/{endpoint:[a-zA-Z0-9/]+}", rpc.NewHandler()).Methods("GET")

	return r
}
