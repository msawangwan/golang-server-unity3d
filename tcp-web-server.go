package main

import (
	"bytes"
	"encoding/binary"
	//"io"
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

/* Decodes an array of bytes that represents a data frame, from Network Byte Order into an associated uint32 value and adds them to a slice to be returned. */
func (dp DataPacket) decodeNBO() []uint32 {
	decodedFrame := make([]uint32, 256)
	var start, end, offset = 0, 4, 0

	for i := 0; i < 1024; i++ { // TODO: get the num iterations based on the size of the stream of data
		if i%4 == 0 {
			var currentVal uint32

			currentSlice := dp.data[start:end]
			buf := bytes.NewReader(currentSlice)
			err := binary.Read(buf, binary.BigEndian, &currentVal)

			if err != nil {
				log.Println("Failed To Decode DataFrame:", err)
			}

			decodedFrame[offset] = currentVal
			offset++

			start += 4
			end += 4
		}
	}

	return decodedFrame
}

func (dp DataPacket) encodeNBO(decodedDataFrame []uint32) []byte {
	encodedFrame := make([]byte, 1024) // use len(dp.data)
	frame := bytes.NewBuffer(encodedFrame)
	//frame := bytes.NewReader(encodedFrame)
	offset := 0

	for i := 0; i < 256; i++ { // TODO: see comment in decoding func
		var currentVal uint32

		currentVal = decodedDataFrame[i]
		//buf := make([]byte, 4)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, currentVal)

		if err != nil {
			log.Println("Failed To Encode DataFrame:", err)
		}

		frame.WriteAt(buf, offset)
		offset += 4
	}

	return frame.Bytes()
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
	log.Println("Terminating Client Connection ...")
}

/* Reads data from a stream once a client establishes a connection. */
func readStream(conn net.Conn) {
	log.Println("Opening Stream For Read/Write ...")

	readChan := make(chan DataPacket)
	errChan := make(chan error)

	go func(readChan chan DataPacket, errChan chan error) {
		for {
			log.Println("Waiting For Client Activity ...")

			data := make([]byte, 1024)
			dataSize, err := conn.Read(data)

			dp := DataPacket{data, dataSize}

			if err != nil {
				errChan <- err
				return
			}

			readChan <- dp
			//go writeStream(conn, dp)
		}
	}(readChan, errChan)

	var isReadComplete bool

	for {
		select { // study selecting channels
		case dp := <-readChan:
			dataFrame := dp.decodeNBO()
			log.Println("Data Buffer Read:", dataFrame)
			reEncoded := dp.encodeNBO(dataFrame)
			go writeStream(conn, reEncoded)
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

func writeStream(conn net.Conn, dataFrame []byte) {
	log.Println("Writing to stream ... ")
	//msg := "This is the server. Only the server can write this reply. If the server is not replying, you will not see this server reply. End of the Message.Te \r\n"
	//log.Println(len(msg))
	//writeChan := make(chan string)
	//writeChan <- msg
	//conn.Write([]byte(msg))
	conn.Write(dataFrame)
}
