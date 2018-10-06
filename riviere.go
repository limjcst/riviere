// BasePath: /rivieve
// swagger:meta
package main

import (
	"github.com/limjcst/riviere/api"
	"github.com/limjcst/riviere/listener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	api.GlobalPool = listener.NewPool()
	api.ServerAddress = ""
	defer api.GlobalPool.Close()
	go func() {
		http.ListenAndServe("127.0.0.1:80", api.NewRouter("/riviere"))
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
