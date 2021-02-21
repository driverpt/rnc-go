package core

import (
	"encoding/binary"
	"io"
)

func ToString(buffer []byte) string {
	return string(buffer)
}

func ToUint8(buffer []byte) uint8 {
	return buffer[0]
}

func ToUInt32BE(buffer []byte) uint32 {
	return binary.BigEndian.Uint32(buffer)
}

func ToUint16BE(buffer []byte) uint16 {
	return binary.BigEndian.Uint16(buffer)
}

func ToUint16LE(buffer []byte) uint16 {
	return binary.LittleEndian.Uint16(buffer)
}

func SkipBytes(reader *io.Reader, length int32) error {
	buffer := make([]byte, length)
	_, err := io.ReadFull(*reader, buffer)

	return err
}

