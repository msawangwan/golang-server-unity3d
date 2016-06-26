package network

import (
	"net"
	"sync"
	"sync/atomic"
)

var idCounter int64 = 0 // placeholder implementation of id

/* Map of current connections by UUID. */
type ConnectedClients map[int]*ClientHandler

/* Manages client connections. */
type ClientController struct {
	Conns         ConnectedClients
	ActiveConnsWG sync.WaitGroup
	sync.Mutex    // use to lock access to the connected clients map
	*ServerLogger
}

/* Create a client controller -- one instances handles all clients. */
func NewClientController() *ClientController {
	return &ClientController{
		Conns:        make(ConnectedClients),
		ServerLogger: NewServerLogger(),
	}
}

/* Must be run as a goroutine, created per client connection. */
func (cc *ClientController) HandleClientConnection(conn net.Conn) {
	cc.LogStatus("Handling new client conn ...")

	clientConnStatus := make(chan error)

	atomic.AddInt64(&idCounter, 1)
	id := atomic.LoadInt64(&idCounter)
	newClient := NewClientHandler(conn, int(id))

	cc.Lock()
	cc.Conns[newClient.UUID] = newClient
	cc.LogStatus("Current number of active conns: ", len(cc.Conns))
	cc.Unlock()

	cc.ActiveConnsWG.Add(1)
	go newClient.Moniter(clientConnStatus)
	// TODO: read channel from client connection on disconnect
	//go func() {
}
