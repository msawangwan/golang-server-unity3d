package main

import (
	"log"
	"net"
	"os"
	"time"
)

const (
	ADDR = ":9080"
)

type DataPacket struct {
	data []byte
	size int
}

func (dp DataPacket) parse() string {
	return string(dp.data[:dp.size])
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	listen()
}

/* Listens for incoming client connections and handles them via goroutines. */
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

/* Moderates conversation with connected client. */
func handleClientStream(conn net.Conn) {
	defer conn.Close()

	readStream(conn)
	conn.Write([]byte("Reply"))
}

/* Reads data from a stream once a client establishes a connection. */
func readStream(conn net.Conn) {
	isOver := false
	timeoutCount := 0

	readChan := make(chan DataPacket)
	errChan := make(chan error)

	go func(readChan chan DataPacket, errChan chan error) {
		for {
			log.Println("Reading Stream ....")

			data := make([]byte, 512)
			dataSize, err := conn.Read(data)

			dp := DataPacket{data, dataSize}

			if err != nil {
				errChan <- err
				return
			}

			readChan <- dp
		}
	}(readChan, errChan)

	ticker := time.Tick(time.Second)

	for {
		select { // study selecting channels
		case dp := <-readChan:
			dataString := dp.parse()
			log.Println("Data Buffer:", dataString)
		case err := <-errChan:
			isOver = true
			log.Println("Closing Stream:", err)

			break
		case <-ticker:
			timeoutCount++
			log.Println("Idle Stream ...", timeoutCount)
		}

		if isOver == true {
			close(readChan)
			close(errChan)

			break
		}
	}
}
