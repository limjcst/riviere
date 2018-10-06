package api

import (
	"github.com/gorilla/mux"
	"github.com/limjcst/riviere/listener"
	"net/http"
)

// GlobalPool stores the current listeners
var GlobalPool *listener.Pool

// NewRouter creates a router
func NewRouter(prefix string) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(prefix+"/spec", GetSpecEndpoint).Methods("GET")
	return router
}

// GetSpecEndpoint returns the swagger doc of apis
// swagger:route GET /spec spec getSpec
//
// Return api spec in swagger format
//
// This will show all available apis.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
func GetSpecEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, "swagger.json")
}
