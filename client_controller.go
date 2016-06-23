package main

import (
	"log"
	"sync"
	"sync/atomic"
)

var idCounter int64 = 0

type ConnectedClients map[int]ClientHandler

type ClientController struct {
	Conns   ConnectedClients
	muConns sync.Mutex
}

func NewClientController() *ClientController {
	return &ClientController{
		Conns:   make(ConnectedClients),
		muConns: &sync.Mutex{},
	}
}

/* Must be run as a goroutine, created per client connection. */
func (cc *ClientController) HandleClientConnection(conn net.Conn) {
	log.Println("handling new client conn")
	defer conn.Close()
	atomic.AddInt64(&idCounter, 1)
	id := atomic.LoadInt64(&idCounter)
	newClient := NewClientConnection(conn, id)

	cc.muConns.Lock()
	cc.Conns[newClient.UUID] = newClient
	cc.muConns.Unlock()

	go newClient.MoniterConnection()
}
