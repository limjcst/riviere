package listener

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// NewServer create a socket listening on a given address
func NewServer(address string, port int) net.Listener {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		// handle error
		ln = nil
	}
	return ln
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
