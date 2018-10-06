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
//
// swagger:operation GET /spec spec getSpec
// ---
// responses:
//   '200':
//     description: api spec
func GetSpecEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, "swagger.json")
}

// TunnelBody is a schema for the api of resource /tunnel
// swagger:parameters addTunnel
type TunnelBody struct {
	// Tunnel parameters
	//
	// required: true
	// in: body
	Body TunnelParam `json:"body"`
}

// TunnelParam is the content of TunnelBody
// swagger:model tunnel
type TunnelParam struct {
	// A port of the gate
	//
	// required: true
	Port int `json:"port"`
	// The address of the target host
	//
	// required: true
	ForwardAddress string `json:"forward_address"`
	// The port of the target host
	//
	// required: true
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
// swagger:operation POST /tunnel tunnel addTunnel
// ---
// responses:
//   '201':
//     description: Add the tunnel successfully.
//   '400':
//     description: Bad request.
//   '409':
//     description: Duplicated post. Port is ocuppied.
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

// PortBody is a schema for the api of resource /tunnel
// swagger:parameters deleteTunnel
type PortBody struct {
	// Port parameters
	//
	// required: true
	// in: body
	Body PortParam `json:"body"`
}

// PortParam is the schema with just a port
// swagger:model
type PortParam struct {
	// A port of the gate
	//
	// required: true
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
// swagger:operation DELETE /tunnel tunnel deleteTunnel
// ---
// responses:
//   '202':
//     description: Delete the tunnel successfully.
//   '400':
//     description: Bad request.
//   '404':
//     description: Port is free.
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
