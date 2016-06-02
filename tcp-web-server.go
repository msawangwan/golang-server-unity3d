package main

import (
	"io"
	"log"
	"net"
)

var ADDR = "127.0.0.1:8080"

func main() {
	listen()
}

func listen() {
	sock, err := net.Listen("tcp", ADDR)

	if err != nil {
		log.Fatal("Error Opening Listener Socket:", err)
	}

	log.Println("Accepting Connections on:", ADDR)

	defer sock.Close()
	for {
		conn, err := sock.Accept()

		if err != nil {
			log.Fatal("Error Accepting Client Connection:", err)
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	buffer := make([]byte, 1500)

	dataSize, err := conn.Read(buffer)

	if err != nil {
		log.Fatal("Error Reading Buffer:", err)
	}

	data := buffer[:dataSize]
	parseData(data)

	conn.Write([]byte("Reply From Server: Data Recieved."))
}

func parseData(data []byte) string {
	return string(data)
}

func receiveVarSizedData(conn net.Conn) []byte {
	totalRecv := 0
	totalLeft := 0

	expectedSize := make([]byte, 4)
	bufferSize, err := conn.Read(expectedSize)

	if err != nil {
		log.Fatal("Error reading buffer size:", err)
	}

	dataBuffer := make([]byte, 0, bufferSize)
	tmpBuffer := make([]byte, 32)

	for totalRecv < bufferSize {
		recv, err := conn.Read(tmpBuffer)

		if err != nil {
			if err != io.EOF {
				log.Fatal("Error receiving buffered data:", err)
			}
			break
		}

		dataBuffer = append(dataBuffer, tmpBuffer[:recv]...)

		totalRecv += recv
		totalLeft -= recv
	}

	return dataBuffer
}
