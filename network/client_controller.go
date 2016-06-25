package network

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var idCounter int64 = 0

/* Map of current connections by UUID. */
type ConnectedClients map[int]*ClientHandler

/* Manages client connections. */
type ClientController struct {
	Conns      ConnectedClients
	sync.Mutex // use to lock access to the connected clients map
}

/* Create a client controller -- one instances handles all clients. */
func NewClientController() *ClientController {
	return &ClientController{
		Conns: make(ConnectedClients),
	}
}

/* Must be run as a goroutine, created per client connection. */
func (cc *ClientController) HandleClientConnection(conn net.Conn) {
	log.Println("handling new client conn")

	atomic.AddInt64(&idCounter, 1)
	id := atomic.LoadInt64(&idCounter)
	newClient := NewClientHandler(conn, int(id))

	cc.Lock()
	cc.Conns[newClient.UUID] = newClient
	log.Println("current num active conns: ", len(cc.Conns))
	cc.Unlock()

	go newClient.Moniter()
}
