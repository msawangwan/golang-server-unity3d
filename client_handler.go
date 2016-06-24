package main

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
func (ch *ClientHandler) Moniter() {
	go ch.stream()
	ch.handleStreamData()
}

func (ch *ClientHandler) stream() {
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
		}

		data := recvBuffer[:bytesRead]
		ch.recvChan <- data
	}
}

func (ch *ClientHandler) handleStreamData() {
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
			break
		}
	}
}
