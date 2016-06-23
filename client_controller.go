package main

import (
	"sync"
	"sync/atomic"
)

var idCounter int64 = 0

type ConnectedClients map[int]ClientHandler

type ClientController struct {
	Conns   ConnectedClients
	muConns sync.Mutex
}

func New() *ClientController {
	return &ClientController{
		Conns:   make(ConnectedClients),
		muConns: &sync.Mutex{},
	}
}

/* Must be run as a goroutine, created per client connection. */
func (cc *ClientController) HandleClientConnection(conn net.Conn) {
	atomic.AddInt64(&idCounter, 1)
	id := atomic.LoadInt64(&idCounter)
	newClient := ClientHandler.New(conn, id)

	cc.muConns.Lock()
	cc.Conns[newClient.UUID] = newClient
	cc.muConns.Unlock()

	go newClient.MoniterConnection()
}
