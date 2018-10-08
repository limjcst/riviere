package listener

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

// CheckPort returns whether a port is listened
func CheckPort(port int) (ok bool) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 1*time.Second)
	if err == nil {
		conn.Close()
		ok = true
	} else {
		ok = false
	}
	return ok
}

func GetServer(t *testing.T) (*TaggedServer, int) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("Cannot listen on any port")
	}
	port := listener.Addr().(*net.TCPAddr).Port
	server := &TaggedServer{
		listener: listener,
		done:     false,
		mutex:    &sync.Mutex{},
	}
	return server, port
}

func StartReadOnlyServer(server *TaggedServer, text string) {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			// handle error
			break
		}
		reader, writer := net.Pipe()
		defer reader.Close()
		defer writer.Close()
		go Response(conn, reader)
		writer.Write([]byte(text))
	}
}

func TestDuplicatedPort(t *testing.T) {
	server1, port := GetServer(t)
	server2 := NewServer("127.0.0.1", port)
	if server2 != nil {
		t.Errorf("Listen on the same port %d", port)
	}
	server1.Close()
	server2 = NewServer("127.0.0.1", port)
	if server2 == nil {
		t.Errorf("Cannot listen on port %d", port)
	}
	server2.Close()
}

func TestResponse(t *testing.T) {
	server, port := GetServer(t)
	text := "Hello World\n"
	go StartReadOnlyServer(server, text)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 1*time.Second)
	if err != nil {
		t.Errorf("Cannot connect to mock server!")
	}
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil || status != text {
		t.Errorf("Response method lost information! %s", status)
	}
}

func TestStart(t *testing.T) {
	server, port := GetServer(t)
	text := "Hello World\n"
	buf := make([]byte, 4096)
	go StartReadOnlyServer(server, text)
	listener, listenerPort := GetServer(t)
	go listener.Start("127.0.0.1", port)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", listenerPort), 1*time.Second)
	if err != nil {
		t.Errorf("Cannot connect to listener server!")
	}
	var status string
	status, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil || status != text {
		t.Errorf("Information lost through listener! %s", status)
	}
	// Test closing target server
	server.Close()
	var n int
	conn.SetReadDeadline(time.Now())
	n, err = conn.Read(buf)
	if err == nil {
		t.Errorf("Receive message after target server closed! %s", buf[:n])
	}
	conn.Close()
	// Test closing listener
	listener.Close()
	if CheckPort(listenerPort) {
		t.Errorf("Connect to listener server after closed!")
	}
}

func TestEcho(t *testing.T) {
	server, port := GetServer(t)
	go func() {
		conn, err := server.listener.Accept()
		if err != nil {
			return
		}
		go Response(conn, conn)
	}()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 1*time.Second)
	if err != nil {
		t.Errorf("Cannot connect to listener!")
	}
	var text string
	buf := make([]byte, 128)
	text = "Hello\n"
	copy(buf, text)
	conn.Write(buf)
	var status string
	status, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil || status != text {
		t.Errorf("Response method lost information! %s", status)
	}
	// Continue
	time.Sleep(1 * time.Second)
	text = "World\n"
	copy(buf, text)
	conn.Write(buf)
	status, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil || status != text {
		t.Errorf("Failed to continue! %s", status)
	}
}
