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

type DataFrameInt struct {
	EncodedData    EncodedFrameInt
	DecodedData    DecodedFrameInt
	PayloadDecoded []uint32
	PayloadEncoded []byte
	PayloadType    PayloadType
	Size           int
}

type EncodedFrameInt struct {
	Payload []byte
	Size
}

type DecodedFrameInt struct {
	Payload []uint32
	Size
}

func (df *DataFrameInt) EncodeIntNetworkByteOrder(unencodedFrame []uint32) error {
	encodedFrame := make([]byte, 0, 1024)

	for i := 0; i < len(unencodedFrame); i++ {
		currentValue := unencodedFrame[i]
		buf := &bytes.Buffer{}

		err := binary.Write(buf, binary.BigEndian, currentValue)
		if err != nil {
			return errors.New("Failed To Encode DataFrame:", err)
		}

		encodedFrame = append(encodedFrame, buf.Bytes()...)
	}

	df.payloadEncoded = encodedFrame
	return nil
}

func (df *DataFrameInt) DecodeIntNetworkByteOrder() error {
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

	df.payloadDecoded = decodedFrame
	return nil
}
