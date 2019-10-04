package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/limjcst/riviere/config"
	"github.com/limjcst/riviere/listener"
	"github.com/limjcst/riviere/models"
	"log"
	"net/http"
)

// ContextInjector injects database into endpoints
type ContextInjector struct {
	db   models.Database
	conf *config.Config
}

// GlobalPool stores the current listeners
var GlobalPool *listener.Pool

// NewRouter creates a router
func NewRouter(c *config.Config, db models.Database) (router *mux.Router) {
	ctx := &ContextInjector{db: db, conf: c}
	router = mux.NewRouter()
	router.HandleFunc(c.Prefix+"/spec", ctx.GetSpecEndpoint).Methods("GET")
	prefix := c.Prefix + "/tunnel"
	router.HandleFunc(prefix, ctx.AddTunnelEndpoint).Methods("POST")
	router.HandleFunc(prefix, ctx.DeleteTunnelEndpoint).Methods("DELETE")
	router.HandleFunc(prefix, ctx.ListTunnelEndpoint).Methods("GET")
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
func (ctx *ContextInjector) GetSpecEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, ctx.conf.Spec)
}

// TunnelBody is a schema for the api of resource /tunnel
// swagger:parameters addTunnel
type TunnelBody struct {
	// Tunnel parameters
	//
	// required: true
	// in: body
	Body models.Tunnel `json:"body"`
}

// NewTunnelParam parses the parameters
func NewTunnelParam(req *http.Request) *models.Tunnel {
	decoder := json.NewDecoder(req.Body)
	param := &models.Tunnel{-1, "", -1}
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
func (ctx *ContextInjector) AddTunnelEndpoint(w http.ResponseWriter, req *http.Request) {
	tunnel := NewTunnelParam(req)
	if tunnel == nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Request to forward %d to %s:%d",
			tunnel.Port, tunnel.ForwardAddress, tunnel.ForwardPort)
		ok := GlobalPool.Listen(tunnel.Port,
			tunnel.ForwardAddress, tunnel.ForwardPort)
		if ok {
			err := ctx.db.NewTunnel(tunnel)
			if err != nil {
				log.Print(err)
			}
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
func (ctx *ContextInjector) DeleteTunnelEndpoint(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	tunnel := PortParam{-1}
	err := decoder.Decode(&tunnel)
	if err != nil || tunnel.Port < 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Request to delete listener on Port %d", tunnel.Port)
		ok := GlobalPool.Delete(tunnel.Port)
		if ok {
			_, err := ctx.db.DeleteTunnel(&models.Tunnel{Port: tunnel.Port})
			if err != nil {
				log.Print(err)
			}
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// ListTunnelEndpoint lists tunnels
// swagger:route GET /tunnel tunnel getTunnel
//
// List tunnels.
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
// swagger:operation GET /tunnel tunnel getTunnel
// ---
// responses:
//   '200':
//     description: Get the tunnels successfully.
//     schema:
//       type: array
//       items:
//         $ref: '#/definitions/Tunnel'
//   '500':
//     description: Internal Server Error.
func (ctx *ContextInjector) ListTunnelEndpoint(w http.ResponseWriter, req *http.Request) {
	var data []byte
	tunnels, err := ctx.db.ListTunnel()
	if err == nil {
		data, err = json.Marshal(tunnels)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
