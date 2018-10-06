package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/limjcst/riviere/listener"
	"log"
	"net/http"
)

// GlobalPool stores the current listeners
var GlobalPool *listener.Pool

// NewRouter creates a router
func NewRouter(prefix string) (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc(prefix+"/spec", GetSpecEndpoint).Methods("GET")
	router.HandleFunc(prefix+"/tunnel", AddTunnelEndpoint).Methods("POST")
	router.HandleFunc(prefix+"/tunnel", DeleteTunnelEndpoint).Methods("DELETE")
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

// TunnelParam is the schema for the api of resource /tunnel
// swagger:parameters addTunnel
type TunnelParam struct {
	// A port of the gate
	//
	// required: true
	// in: body
	Port int `json:"port"`
	// The address of the target host
	//
	// required: true
	// in: body
	ForwardAddress string `json:"forward_address"`
	// The port of the target host
	//
	// required: true
	// in: body
	ForwardPort int `json:"forward_port"`
}

// NewTunnelParam parses the parameters
func NewTunnelParam(req *http.Request) *TunnelParam {
	decoder := json.NewDecoder(req.Body)
	param := &TunnelParam{-1, "", -1}
	err := decoder.Decode(param)
	if err != nil || param.Port < 0 || param.ForwardAddress == "" ||
		param.ForwardPort < 0 {
		param = nil
	}
	return param
}

// AddTunnelEndpoint adds a tunnel
// swagger:route POST /tunnel tunnel addTunnel
//
// Add a tunnel.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       201:
//       400:
//       409:
func AddTunnelEndpoint(w http.ResponseWriter, req *http.Request) {
	tunnel := NewTunnelParam(req)
	if tunnel == nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Request to forward %d to %s:%d",
			tunnel.Port, tunnel.ForwardAddress, tunnel.ForwardPort)
		ok := GlobalPool.Listen(tunnel.Port,
			tunnel.ForwardAddress, tunnel.ForwardPort)
		if ok {
			// TODO: Store the tunnel in the database
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}

// PortParam is the schema with just a port
// swagger:parameters deleteTunnel
type PortParam struct {
	// A port of the gate
	//
	// required: true
	// in: body
	Port int `json:"port"`
}

// DeleteTunnelEndpoint deletes a tunnel
// swagger:route DELETE /tunnel tunnel deleteTunnel
//
// Delete a tunnel.
//
//     Consumes:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       202:
//       400:
//       404:
func DeleteTunnelEndpoint(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	tunnel := PortParam{-1}
	err := decoder.Decode(&tunnel)
	if err != nil || tunnel.Port < 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Request to delete listener on Port %d", tunnel.Port)
		ok := GlobalPool.Delete(tunnel.Port)
		if ok {
			// TODO: Remove the tunnel from the database
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}