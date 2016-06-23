package main

import (
	"log"
	"net"
)

type ClientHandler struct {
	connection net.Conn
	UUID       int
	sendCh     chan []byte
	recvCh     chan []byte
}

func New(conn net.Conn, uuid int) *ClientHandler {
	return &ClientHandler{
		connection: conn,
		UUID:       uuid,
		sendCh:     make(chan []byte),
		recvCh:     make(chan []byte),
	}
}

func (ch *ClientHandler) MoniterConnection() {
	for {
		select {
		case read := <-ch.sendCh:
			log.Println("reading", string(read))
		case write := <-ch.writeCh:
			log.Println("writing", string(write))
		}
	}
}
