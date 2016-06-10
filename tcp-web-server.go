package main

import (
	"log"
	"net"
	"os"
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

	log.Println("Accepting Connections On:", ADDR)

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
	log.Println("Terminating Client Connection")
}

/* Reads data from a stream once a client establishes a connection. */
func readStream(conn net.Conn) {
	log.Println("Opening Stream For Read ... ")

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
			go writeStream(conn)
		}
	}(readChan, errChan)

	var isReadComplete bool

	for {
		select { // study selecting channels
		case dp := <-readChan:
			dataString := dp.parse()
			log.Println("Data Buffer Read:", dataString)
		case err := <-errChan:
			isReadComplete = true
			log.Println("Closing Read Stream:", err)

			break
		}

		if isReadComplete == true {
			close(readChan)
			close(errChan)

			break
		}
	}
}

var hasSentOneMsg = false

func writeStream(conn net.Conn) {
	if hasSentOneMsg == false {
		log.Println("Writing to stream ... ")
		msg := "This is the server. Only the server can write this reply. If the server is not replying, you will not see this server reply. End of the Message.Te \r\n"
		log.Println(len(msg))
		//writeChan := make(chan string)
		//writeChan <- msg
		conn.Write([]byte(msg))
		conn.Write([]byte(msg))
		//hasSentOneMsg = true
	}
}
