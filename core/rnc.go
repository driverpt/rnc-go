package core

import (
	"io"
)

const (
	RNCSignature = "RNC"
	HeaderSizeInBytes = int32(18)
)

type RNCHeader struct {
	Signature string
	CompressionMethod uint8
	OriginalSize uint32
	CompressedSize uint32
	OriginalCRC uint16
	CompressedCRC uint16
	Leeway uint8
	PackChunks uint8
}

type RNCReader interface {
	Unpack() ([]byte, error)
}

func ReadHeader(src io.Reader) (*RNCHeader, error) {
	result := RNCHeader {
	}

	buffer, err := readHeaderBytes(&src, 3)
	if err != nil {
		return nil, NonRNCStream{}
	}

	result.Signature = ToString(buffer)

	if result.Signature != RNCSignature {
		return nil, NonRNCStream{}
	}

	buffer, err = readHeaderBytes(&src, 1)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.CompressionMethod = ToUint8(buffer)

	buffer, err = readHeaderBytes(&src, 4)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.OriginalSize = ToUInt32BE(buffer)

	buffer, err = readHeaderBytes(&src, 4)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.CompressedSize = ToUInt32BE(buffer)

	buffer, err = readHeaderBytes(&src, 2)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.OriginalCRC = ToUint16BE(buffer)

	buffer, err = readHeaderBytes(&src, 2)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.CompressedCRC = ToUint16BE(buffer)

	buffer, err = readHeaderBytes(&src, 1)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.Leeway = ToUint8(buffer)

	buffer, err = readHeaderBytes(&src, 1)
	if err != nil {
		return nil, NonRNCStream{}
	}
	result.PackChunks = ToUint8(buffer)

	return &result, nil
}

func readHeaderBytes(src *io.Reader, size int) ([]byte, error) {
	buffer := make([]byte, size)

	_, err := (*src).Read(buffer)

	if err != nil {
		return nil, NonRNCStream{}
	}
	return buffer, nil
}


