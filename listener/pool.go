package listener

import (
	"log"
)

// Pool stores a list of listeners
type Pool struct {
	listeners map[int]*TaggedServer
}

// NewPool initializes a Pool
func NewPool() *Pool {
	return &Pool{
		listeners: make(map[int]*TaggedServer),
	}
}

// Add a listener to a pool
// return whether succeeded
func (pool *Pool) Add(port int, listener *TaggedServer) bool {
	_, ok := pool.listeners[port]
	if !ok {
		pool.listeners[port] = listener
	}
	return !ok
}

// Delete a listener from a pool and close it
func (pool *Pool) Delete(port int) {
	listener, ok := pool.listeners[port]
	if ok {
		listener.Close()
		delete(pool.listeners, port)
	}
}

// Length returns the size of the pool
func (pool *Pool) Length() int {
	return len(pool.listeners)
}

// Close every listener
func (pool *Pool) Close() {
	for port, listener := range pool.listeners {
		listener.Close()
		log.Printf("Release Port %d", port)
	}
}
