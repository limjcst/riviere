package main

import (
	"fmt"
	"github.com/limjcst/riviere/listener"
	"io"
	"net/http"
	"time"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	})
	http.ListenAndServe("127.0.0.1:80", nil)
}

func main() {
	go server()
	ln := listener.NewServer("127.0.0.1", 8000)
	if ln == nil {
		fmt.Println("port not available")
	} else {
		go func() { ln.Start("127.0.0.1", 80) }()
		time.Sleep(10 * time.Second)
		ln.Close()
	}
}
