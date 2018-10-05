package listener

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// TaggedServer can record whether it is done
type TaggedServer struct {
	listener net.Listener
	done     chan bool
}

// Close a listener and set a tag
func (server *TaggedServer) Close() {
	go func() {
		server.done <- true
	}()
	server.listener.Close()
}

// NewServer create a socket listening on a local address and forward
func NewServer(address string, port int) (server *TaggedServer) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		// handle error
		server = nil
	} else {
		server = &TaggedServer{
			listener: ln,
			done:     make(chan bool),
		}
	}
	return server
}

// Start a server forwarding to given address
func (server *TaggedServer) Start(forwardAddress string, forwardPort int) {
	for {
		conn, err := server.listener.Accept()
		select {
		case <-server.done:
			defer func() {
				if recover() != nil {
				}
			}()
			conn.Close()
			return
		default:
		}
		if err != nil {
			// handle error
			continue
		}
		address := fmt.Sprintf("%s:%d", forwardAddress, forwardPort)
		go func() {
			forwardServer, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				return
			}
			defer forwardServer.Close()
			Response(conn, forwardServer)
		}()
	}
}

// Response to the request of server
func Response(conn net.Conn, provider io.ReadWriter) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(conn, provider)
	}()
	go func() {
		defer wg.Done()
		io.Copy(provider, conn)
	}()
	wg.Wait()
	defer conn.Close()
}
