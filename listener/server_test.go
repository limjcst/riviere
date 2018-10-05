package listener

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"
)

func GetServer(t *testing.T) (net.Listener, int) {
	server, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("Cannot listen on any port")
	}
	port := server.Addr().(*net.TCPAddr).Port
	return server, port
}

func TestDuplicatedPort(t *testing.T) {
	server1, port := GetServer(t)
	server2 := NewServer("127.0.0.1", port)
	if server2 != nil {
		t.Errorf("Listen on the same port %d", port)
	}
	server1.Close()
}

func TestResponse(t *testing.T) {
	server, port := GetServer(t)
	text := "Hello World\n"
	go func(t *testing.T) {
		for {
			conn, err := server.Accept()
			if err != nil {
				// handle error
			}
			go Response(conn, bytes.NewBufferString(text))
		}
	}(t)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Errorf("Cannot connect to mock server!")
	}
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil || status != text {
		t.Errorf("Response method lost information! %s", status)
	}
}