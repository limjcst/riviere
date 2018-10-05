package listener

import (
	"testing"
)

func TestEditPool(t *testing.T) {
	pool := NewPool()
	if pool.Length() != 0 {
		t.Errorf("New Pool is not empty!")
	}
	server, port := GetServer(t)
	ok := pool.Add(port, server)
	len := pool.Length()
	if !ok || len != 1 {
		t.Errorf("Failed to add a listener! Length: %d", len)
	}
	if pool.Add(port, server) {
		t.Errorf("Add duplicated listener")
	}
	pool.Delete(port)
	if pool.Length() != 0 {
		t.Errorf("Failed to delete a listener!")
	}
}

func TestClose(t *testing.T) {
	pool := NewPool()
	server, port := GetServer(t)
	pool.Add(port, server)
	if !CheckPort(port) {
		t.Errorf("Failed to listen port!")
	}
	pool.Close()
	if CheckPort(port) {
		t.Errorf("Failed to close listener with port %d", port)
	}
}
