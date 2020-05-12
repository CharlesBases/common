package middleware

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	router *mux.Router
)

func init() {
	router = mux.NewRouter()
}

func Run() {
	n := negroni.New()

	n.Use(Recovery())
	n.Use(NegroniLogger())
	n.UseFunc(negroni.HandlerFunc(Cors()))
	n.UseFunc(negroni.HandlerFunc(JWT()))

	n.UseHandler(router)

	n.Run(":8080")
}
