package main

import (
	"log"
	"net"
)

type ServerCore struct {
	HostAddr       string
	ListenerSocket net.Listener
	ClientCtrller  ClientController
}

func New(addr string) *ServerCore {
	return &ServerCore{
		HostAddr: addr,
	}
}

func (server *ServerCore) Start() {
	server.bind()
	go server.run()
}

func (server *ServerCore) bind() {
	lnSock, err := net.Listen("tcp", server.HostAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	server.ListenerSocket = lnSock
}

func (server *ServerCore) run() {
	log.Println("Server running ...")
	for {
		conn, err := server.ListenerSocket.Accept()
		if err != nil {
			log.Fatal("%v", err)
		}

		go server.ClientCtrller.HandleClientConnection(conn)
	}
}
