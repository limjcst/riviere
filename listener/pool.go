package listener

import (
	"log"
	"sync"
)

// Pool stores a list of listeners of given address
type Pool struct {
	address   string
	listeners map[int]*TaggedServer
	mutex     *sync.Mutex
}

// NewPool initializes a Pool
func NewPool(address string) *Pool {
	return &Pool{
		address:   address,
		listeners: make(map[int]*TaggedServer),
		mutex:     &sync.Mutex{},
	}
}

// Add a listener to a pool
// return whether succeeded
func (pool *Pool) Add(port int, listener *TaggedServer) bool {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	_, ok := pool.listeners[port]
	if !ok {
		pool.listeners[port] = listener
	}
	return !ok
}

// Delete a listener from a pool and close it
func (pool *Pool) Delete(port int) bool {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	listener, ok := pool.listeners[port]
	if ok {
		listener.Close()
		delete(pool.listeners, port)
	}
	return ok
}

// Length returns the size of the pool
func (pool *Pool) Length() int {
	return len(pool.listeners)
}

// Close every listener
func (pool *Pool) Close() {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	for port, listener := range pool.listeners {
		listener.Close()
		log.Printf("Release Port %d", port)
	}
}

// Listen setups a tunnel between given addresses
func (pool *Pool) Listen(port int,
	forwardAddress string, forwardPort int) (ok bool) {
	ln := NewServer(pool.address, port)
	if ln == nil {
		log.Printf("Port %d is not available", port)
		ok = false
	} else {
		go func() {
			ln.Start(forwardAddress, forwardPort)
		}()
		pool.Add(port, ln)
		ok = true
	}
	return ok
}
