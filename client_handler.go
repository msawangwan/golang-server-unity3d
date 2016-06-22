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

func (c *ClientHandler) Recv() {
	log.Println("Waiting For Request.")

	recvCh := make(chan DataFrame)
	errCh := make(chan error)

	go func(recvCh chan DataFrame, errCh chan error) {
		for {
			recvBuffer := make([]byte, 1024)

			bytesRead, err := c.connection.Read(recvBuffer)
			if err != nil {
				errCh <- err
				return
			}

			data := recvBuffer[:bytesRead]
			df := DataFrame{
				payload: data,
				size: len(data),
			}

			readCh <- df
		}(readCh, errCh)
	}

		var hasDisconnected bool

		for {
			select {
			case df := <-readCh:
				dataFrame, err := df.DecodeNetworkByteOrder()
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
}
