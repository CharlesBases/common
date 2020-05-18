package middleware

import (
	"fmt"
	"net/http"
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
	n.Use(Negroni())
	n.Use(Cors())
	n.Use(JWT())

	n.UseHandler(router)
	router.HandleFunc("/", Holle)

	n.Run(":8080")
}

func Holle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "this is home")
}
