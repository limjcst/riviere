// Package classification Rivière
//
// Set up tunnels between local ports and remote addresses dynamically.
//
//     BasePath: /rivieve
//     Version: Beta
// swagger:meta
//go:generate swagger generate spec -o swagger.json
package main

import (
	"github.com/gorilla/handlers"
	"github.com/limjcst/riviere/api"
	"github.com/limjcst/riviere/listener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Manage ports of each address available
	api.GlobalPool = listener.NewPool("")
	defer api.GlobalPool.Close()
	go func() {
		http.ListenAndServe("127.0.0.1:80",
			handlers.CORS(
				handlers.AllowedOrigins([]string{"*"}))(
				api.NewRouter("/riviere")))
	}()
	log.Printf("Rivière has started")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	log.Printf("Bye Bye!")
}
