package main

import (
	"log"
	"net"
)

var ADDR = "127.0.0.1:8080"

func main() {
}

func listen() {
	sock, err := net.Listen("tcp", ADDR)

	if err != nil {
		log.Fatal("Error Opening Listener Socket:", err)
	}

	defer sock.Close()
	for {
		conn, err := sock.Accept()

		if err != null {
			log.Fatal("Error Accepting Client Connection:", err)
			return
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

	conn.Write([]byte("Data Recieved."))
}

func parseData(data []byte) {
	dataString := string(data)
	log.Println(dataString)
}
