package main

import (
	"bytes"
	"encoding/binary"
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
	decodedFrame := make([]uint32, 0, 256)
	var start, end = 0, 4

	for i := 0; i < dp.size; i++ {
		if i%4 == 0 {
			var currentVal uint32

			currentSlice := dp.data[start:end]
			buf := bytes.NewReader(currentSlice)
			err := binary.Read(buf, binary.BigEndian, &currentVal)

			if err != nil {
				log.Println("Failed To Decode DataFrame:", err)
			}

			decodedFrame = append(decodedFrame, currentVal)

			start += 4
			end += 4
		}
	}

	return decodedFrame
}

/* Encodes a slice of decoded data into a slice of bytes, in network byte order. */
func (dp DataPacket) encodeNBO(decodedFrame []uint32) []byte {
	encodedFrame := make([]byte, 0, 1024)

	for i := 0; i < len(decodedFrame); i++ {
		currentVal := decodedFrame[i]

		buf := &bytes.Buffer{} // OR: buf := new(bytes.buffer)

		err := binary.Write(buf, binary.BigEndian, currentVal)
		if err != nil {
			log.Println("Failed To Encode DataFrame:", err)
		}

		encodedFrame = append(encodedFrame, buf.Bytes()...)
	}

	return encodedFrame
}

func init() {
	log.SetOutput(os.Stdout)
}

/* Listens for incoming client connections and handles them via goroutines. */
func Listen() {
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

			recvBuffer := make([]byte, 1024)
			bytesRead, err := conn.Read(recvBuffer)

			data := recvBuffer[:bytesRead]
			dp := DataPacket{data, len(data)}

			if err != nil {
				errChan <- err
				return
			}

			readChan <- dp
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
	frameSize := len(dataFrame)
	log.Println("Size Of DataFrame:", frameSize)

	encodedSize := make([]byte, 4)
	binary.BigEndian.PutUint32(encodedSize, uint32(frameSize))

	conn.Write(encodedSize)
	conn.Write(dataFrame)
}
