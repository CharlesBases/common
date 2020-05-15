package middleware

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	router *mux.Router
)

func init() {
	router = mux.NewRouter()
}

func TestMiddleware(t *testing.T) {
	n := negroni.New()

	n.Use(Recovery())
	n.UseFunc(negroni.HandlerFunc(Negroni()))
	n.UseFunc(negroni.HandlerFunc(Cors()))
	n.UseFunc(negroni.HandlerFunc(JWT()))

	n.UseHandler(router)

	n.Run(":8080")
}
