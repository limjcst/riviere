package main

import (
	"fmt"
	"github.com/limjcst/riviere/listener"
	"io"
	"net"
	"net/http"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	})
	http.ListenAndServe(":80", nil)
}

func main() {
	go server()
	helloServer, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		fmt.Println("Cannot connect to mock server!")
	}
	ln := listener.NewServer("127.0.0.1", 8000)
	if ln == nil {
		fmt.Println("port not available")
	} else {
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
			}
			go listener.Response(conn, helloServer)
		}
	}
}
