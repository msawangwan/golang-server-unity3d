package main

import (
	"io"
	"log"
	"net"
	"time"
)

var ADDR = ":9080"

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

		go handleClientStream(conn)
	}
}

func handleClientStream(conn net.Conn) {
	readStream(conn)
	conn.Write([]byte("Reply"))
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

var timeoutCount = 0

func readStream(conn net.Conn) {
	ch := make(chan []byte)
	chErr := make(chan error)

	go func(ch chan []byte, chErr chan error) {
		for {
			data := make([]byte, 512)
			_, err := conn.Read(data)

			if err != nil {
				chErr <- err
				return
			}

			ch <- data
		}
	}(ch, chErr) // <- to research

	ticker := time.Tick(time.Second)

	isOver := false
	for {
		select {
		case data := <-ch:
			log.Println(parseData(data))
		case err := <-chErr:
			log.Println(err)

			if err == io.EOF {
				isOver = true
			}

			break
		case <-ticker:
			timeoutCount++
			log.Println("timeout ...", timeoutCount)
		}

		if isOver == true {
			timeoutCount = 0
			close(ch)
			close(chErr)
			break
		}
	}
}
