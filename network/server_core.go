package network

import (
	"net"
	"sync"
	"time"
)

/* The server module, the central unit. */
type ServerCore struct {
	HostAddr       string // switch to unexported fields??
	ClientCtr      *ClientController
	ListenerSocket net.Listener
	ListenerWG     sync.WaitGroup
	killSignal     chan bool
	*ServerLogger
}

func NewServerInstance(addr string) *ServerCore {
	server := &ServerCore{
		HostAddr:     addr,
		ClientCtr:    NewClientController(),
		killSignal:   make(chan bool, 1),
		ServerLogger: NewServerLogger(),
	}

	server.ListenerWG.Add(1)
	return server
}

func (server *ServerCore) Start() {
	server.setup()
	server.bind()
	go server.run()
}

func (server *ServerCore) Shutdown(shutdownNow bool) {
	if shutdownNow {
		server.kill()
	} else {
		defer server.ListenerSocket.Close()
	}

	server.ClientCtr.ActiveConnsWG.Wait() // TODO: need to call done for this wg
	server.ListenerWG.Wait()
}

func (server *ServerCore) setup() {
	server.LogStatus("Server core starting ...")
}

func (server *ServerCore) bind() {
	lnSock, err := net.Listen("tcp", server.HostAddr)
	if err != nil {
		server.LogFatalAlert("Error occur binding to socket: ", err)
	}

	server.ListenerSocket = lnSock
}

func (server *ServerCore) run() {
	server.LogStatus("Server core running ...")
	defer server.ListenerWG.Done()

	for {
		select {
		case <-server.killSignal:
			server.LogStatus("Kill signal received")
			server.ListenerSocket.Close()
			close(server.killSignal)
			return
		default:
			server.LogStatus("Accepting connections ... ")
		}

		if tcpListener, ok := server.ListenerSocket.(*net.TCPListener); ok { // get the base type Conn
			tcpListener.SetDeadline(time.Now().Add(1000 * time.Millisecond))
		}

		conn, err := server.ListenerSocket.Accept()
		if err != nil {
			if timeout, ok := err.(*net.OpError); ok && timeout.Timeout() {
				continue
			}
			server.LogFatalAlert("Error during client connection attempt", err)
		}

		server.LogStatus("Incoming client connection request ...")
		go server.ClientCtr.HandleClientConnection(conn)
	}
}

// TODO: Close the channel on shutdown.
func (server *ServerCore) kill() {
	server.killSignal <- true
}
