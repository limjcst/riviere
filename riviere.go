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
	"flag"
	"fmt"
	_ "github.com/go-swagger/go-swagger"
	"github.com/gorilla/handlers"
	"github.com/limjcst/riviere/api"
	"github.com/limjcst/riviere/listener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Name is the name of this application
const Name = "Rivière"

func start(host *string, port *int, configFile *string) {
	// Manage ports of each address available
	api.GlobalPool = listener.NewPool("")
	defer api.GlobalPool.Close()
	address := fmt.Sprintf("%s:%d", *host, *port)
	var c config
	c.parseConfig(*configFile)
	router := api.NewRouter(c.Prefix, c.DBDriver, c.DBSourceName)
	if router == nil {
		return
	}
	go func() {
		http.ListenAndServe(address,
			handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router))
	}()
	log.Printf("%s has started: %s", Name, address)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	log.Printf("Bye Bye!")
}

func main() {
	flagSet := flag.NewFlagSet(Name, flag.ContinueOnError)
	host := flagSet.String("host", "127.0.0.1",
		"Host address. It's dangerous to be not localhost")
	port := flagSet.Int("port", 80, "Port")
	configFile := flagSet.String("config", "config.yml", "Config file path.")
	err := flagSet.Parse(os.Args[1:])
	if err != flag.ErrHelp {
		start(host, port, configFile)
	}
}
