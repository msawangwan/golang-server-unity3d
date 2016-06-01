package main

import (
	"log"
	"net"
)

var ADDR = ":8080"

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

func parseData(data []byte) {
	dataString := string(data)
	log.Println(dataString)
}
