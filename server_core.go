package main

import (
	"log"
	"net"
	"sync"
)

/* The server module, the central unit. */
type ServerCore struct {
	HostAddr       string
	ClientCtr      *ClientController
	ListenerSocket net.Listener // consider using unexported fields
	ListenerWG     sync.WaitGroup
}

func NewServerInstance(addr string) *ServerCore {
	return &ServerCore{
		HostAddr:  addr,
		ClientCtr: NewClientController(),
	}
}

func (server *ServerCore) Start() {
	server.bind()

	server.ListenerWG.Add(1)
	go server.run()
}

func (server *ServerCore) Shutdown() {
	defer server.ListenerSocket.Close()
	server.ListenerWG.Wait()
}

func (server *ServerCore) bind() {
	lnSock, err := net.Listen("tcp", server.HostAddr)
	if err != nil {
		log.Fatalf("Error on server bind to socket: ", err)
	}

	server.ListenerSocket = lnSock
}

func (server *ServerCore) run() {
	log.Println("Server running ...")
	defer server.ListenerWG.Done()

	for {
		conn, err := server.ListenerSocket.Accept()
		if err != nil {
			log.Fatal("Error on accepting client connection: ", err)
		}
		log.Println("new client connecting ...")
		go server.ClientCtr.HandleClientConnection(conn)
	}
}

/* turn this into a test func in a server-core_test.go file */
func (server *ServerCore) handleClientConn(conn net.Conn) {
	log.Println("new client conn")
}
