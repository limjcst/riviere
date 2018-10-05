package main

import (
	"fmt"
	"github.com/limjcst/riviere/listener"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	})
	http.ListenAndServe("127.0.0.1:80", nil)
}

func main() {
	go server()
	pool := listener.NewPool()
	defer pool.Close()
	ln := listener.NewServer("127.0.0.1", 8000)
	if ln == nil {
		fmt.Println("port not available")
	} else {
		go func() { ln.Start("127.0.0.1", 80) }()
		pool.Add(8000, ln)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	log.Printf("Bye Bye!")
}
