package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type PayloadType uint8

const (
	INT_TYPE PayloadType = iota
	STRING_TYPE
)

type DataFrame struct {
	payload     []byte
	payloadType PayloadType
	size        int
}

func (df *DataFrame) DecodeNetworkByteOrder() ([]uint32, error) {
	decodedFrame := make([]uint32, 0, 256)
	var start, end = 0, 4

	for i := 0; i < df.size; i++ {
		if i%4 == 0 {
			var currentValue uint32

			currentSlice := df.payload[start:end]
			tempBuffer := bytes.NewReader(currentSlice)

			err := binary.Read(buf, binary.BigEndian, &currentValue)
			if err != nil {
				return nil, errors.New("Error Decoding Dataframe:", err)
			}

			decodedFrame = append(decodedFrame, currentValue)

			start += 4
			end += 4
		}
	}

	return decodedFrame, nil
}
