package main

import (
	"log"
	"net"
)

type ClientHandler struct {
	connection net.Conn
	UUID       int
	sessionID  int
}

func New(conn net.Conn, uuid int, sessionid int) *ClientHandler {
	return &ClientHandler{
		connection: conn,
		UUID:       uuid,
		sessionID:  sessionid,
	}
}

func (c *ClientHandler) Recv() {
	log.Println("Waiting For Request.")

	recvCh := make(chan DataFrameInt)
	errCh := make(chan error)

	go func(recvCh chan DataFrameInt, errCh chan error) {
		for {
			recvBuffer := make([]byte, 1024)

			bytesRead, err := c.connection.Read(recvBuffer)
			if err != nil {
				errCh <- err
				return
			}

			data := recvBuffer[:bytesRead]
			df := DataFrameInt{
				payloadEncoded: data,
				size:           len(data),
			}

			readCh <- df
		}
	}(readCh, errCh)

	var hasDisconnected bool

	for {
		select {
		case df := <-readCh:
			err := df.DecodeIntNetworkByteOrder()
			if err != nil {
				log.Println(err)
			}
		case err := <-errCh:
			hasDisconnected = true
			log.Println(err)
			break
		}

		if hasDisconnected == true {
			close(readCh)
			close(errCh)
			break
		}
	}
}

func (c *ClientHandler) Send(dataFrame *DataFrameInt) {
	frameSize := dataFrame.Size
}
