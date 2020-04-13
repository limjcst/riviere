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
	done     bool
	mutex    *sync.Mutex
}

// Close a listener and set a tag
func (server *TaggedServer) Close() {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	server.done = true
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
			done:     false,
			mutex:    &sync.Mutex{},
		}
	}
	return server
}

// Start a server forwarding to given address
func (server *TaggedServer) Start(forwardAddress string, forwardPort int) {
	for {
		server.mutex.Lock()
		done := server.done
		server.mutex.Unlock()
		if done {
			break
		}
		conn, err := server.listener.Accept()
		if err != nil {
			// handle error
			continue
		}
		address := fmt.Sprintf("%s:%d", forwardAddress, forwardPort)
		go func() {
			forwardServer, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				conn.Close()
				return
			}
			Response(conn, forwardServer)
		}()
	}
}

// Response to the request of server
func Response(conn net.Conn, provider net.Conn) {
	ch := make(chan bool)
	go transfer(conn, provider, ch)
	go transfer(provider, conn, ch)
	<-ch
	conn.Close()
	provider.Close()
	<-ch
	close(ch)
}

func transfer(receiver io.Writer, provider io.Reader, ch chan bool) {
	io.Copy(receiver, provider)
	ch <- true
}
