package network

import (
	"io"
	"log"
	"net"
)

const (
	RECV_BUFFER_MAX_SIZE = 1024 // read buff max size
)

type ClientHandler struct {
	connection       net.Conn
	UUID             int
	sendChan         chan []byte
	recvChan         chan []byte
	disconnectedChan chan error // a value on this channel signals a dced client
}

func NewClientHandler(conn net.Conn, uuid int) *ClientHandler {
	return &ClientHandler{
		connection:       conn,
		UUID:             uuid,
		sendChan:         make(chan []byte),
		recvChan:         make(chan []byte),
		disconnectedChan: make(chan error),
	}
}

/* Moniter the client conn stream as it sends/recvs data (or disconnects). */
func (ch *ClientHandler) Moniter(connStatus chan error) {
	go ch.beginReadStream()
	defer close(connStatus) // closing the chan sends a disconn signal
	ch.handleStreamData()
	//status := ch.handleStreamData()
	//connStatus <- status
}

func (ch *ClientHandler) beginReadStream() {
	for {
		log.Println("stream ready")
		recvBuffer := make([]byte, RECV_BUFFER_MAX_SIZE)

		bytesRead, err := ch.connection.Read(recvBuffer)
		if err != nil {
			if err == io.EOF {
				ch.disconnectedChan <- err
				return
			}
			log.Println("error on recv", err)
			continue
		}

		data := recvBuffer[:bytesRead]
		ch.recvChan <- data
	}
}

/* TODO: change from []byte to a struct with fields lenPrefix and encoded */
func (ch *ClientHandler) beginWriteStream(encodedData []byte) {
	lenPrefix := len(encodedData)
	ch.connection.Write(lenPrefix)
	ch.connection.Write(encodedData)
}

/* Detects read/write ops on a stream - returns an error that represents a
disconnect signal when the client has closed the connection. */
func (ch *ClientHandler) handleStreamData() error {
	defer ch.connection.Close()
	defer close(ch.sendChan)
	defer close(ch.recvChan)
	defer close(ch.disconnectedChan)

	for {
		select {
		case send := <-ch.sendChan:
			log.Println("sending", send)
		case recv := <-ch.recvChan:
			log.Println("recving", recv)
		case disconnect := <-ch.disconnectedChan:
			log.Println("client disconnected", disconnect)
			return disconnect
		}
	}
}
