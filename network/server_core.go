package network

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
	*ServerLogger
}

func NewServerInstance(addr string) *ServerCore {
	return &ServerCore{
		HostAddr:     addr,
		ClientCtr:    NewClientController(),
		ServerLogger: NewServerLogger(),
	}
}

func (server *ServerCore) Start() {
	server.setup()
	server.bind()

	server.ListenerWG.Add(1)
	go server.run()
}

func (server *ServerCore) Shutdown() {
	defer server.ListenerSocket.Close()
	server.ListenerWG.Wait()
}

func (server *ServerCore) setup() {
	server.LogInfo("Server core starting ...")
}

func (server *ServerCore) bind() {
	lnSock, err := net.Listen("tcp", server.HostAddr)
	if err != nil {
		server.LogFatalAlert("Error on server bind to socket: ", err)
	}

	server.ListenerSocket = lnSock
}

func (server *ServerCore) run() {
	server.LogInfo("Server core running ...")
	defer server.ListenerWG.Done()

	for {
		conn, err := server.ListenerSocket.Accept()
		if err != nil {
			server.LogFatalAlert("Error on accepting client connection: ", err)
		}
		server.LogInfo("new client connecting ...")
		go server.ClientCtr.HandleClientConnection(conn)
	}
}

/* turn this into a test func in a server-core_test.go file */
func (server *ServerCore) handleClientConn(conn net.Conn) {
	log.Println("new client conn")
}
