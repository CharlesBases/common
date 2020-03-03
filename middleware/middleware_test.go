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

	n.Use(NegroniLogger())
	n.Use(Recovery())
	n.UseFunc(negroni.HandlerFunc(Cors()))

	n.UseHandler(router)

	n.Run(":8080")
}
