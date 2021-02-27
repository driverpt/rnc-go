package rnc1

import (
	"github.com/driverpt/rnc-go/core"
	"io"
)

const DefaultBufferSize = 1024

type BitStream struct {
	bitBuffer      uint32
	bitBufferCount int
	endOfData      int64
	bytesRead      int64
	reader         *io.Reader

	reachedEOF       bool
	bufferSize       int32
	byteBuffer       []byte
	currentByteIndex int32
	consumedBytes    int
	lastReadBytes    []byte
}

func NewBitStream(reader io.Reader, initialPos int32, endPos int32) BitStream {
	result := BitStream{
		reader:         &reader,
		bitBufferCount: 16,
		bufferSize:     DefaultBufferSize,
		bytesRead:      int64(initialPos),
		endOfData:      int64(endPos),
	}

	result.byteBuffer = make([]byte, result.bufferSize)

	_, err := result.refreshByteBuffer()
	if err != nil {
		panic(err)
	}

	if initialPos != 0 {
		result.advanceByteBufferIndex(int32(result.bytesRead))
	}

	result.bitBuffer = result.readU16LE()

	return result
}

func (s *BitStream) RefreshBuffer() {
	s.bitBufferCount -= 16
	s.bitBuffer &= (1 << uint32(s.bitBufferCount)) - 1

	s.refreshByteBuffer()

	if s.bytesRead < s.endOfData-1 {
		s.bitBuffer |= s.readU16LE() << s.bitBufferCount
		s.bitBufferCount += 16
	} else if s.bytesRead == s.endOfData-1 {
		s.bitBuffer |= s.readByte() << s.bitBufferCount
		s.bitBufferCount += 16
	}
}

func (s *BitStream) Advance(bits int) {
	s.bitBuffer >>= bits
	s.bitBufferCount -= bits

	if s.bitBufferCount >= 16 {
		return
	}

	s.advanceByteBufferIndex(2)

	if s.bytesRead < s.endOfData-1 {
		s.bitBuffer |= s.readU16LE() << s.bitBufferCount
		s.bitBufferCount += 16
	} else if s.bytesRead < s.endOfData {
		s.bitBuffer |= s.readByte() << s.bitBufferCount
		s.bitBufferCount += 16
	}
}

func (s *BitStream) advanceByteBufferIndex(length int32) {
	s.currentByteIndex += length
	s.bytesRead += int64(length)
	// Ensure that there's at least 2 contiguous bytes in buffer available to read
	if s.currentByteIndex+2 > s.bufferSize {
		s.refreshByteBuffer()
	}
}

func (s *BitStream) readU16LE() uint32 {
	result := s.byteBuffer[s.currentByteIndex : s.currentByteIndex+2]
	return uint32(core.ToUint16LE(result))
}

func (s *BitStream) readByte() uint32 {
	return uint32(s.byteBuffer[s.currentByteIndex])
}

func (s *BitStream) Peek(mask uint32) uint32 {
	return s.bitBuffer & mask
}

func (s *BitStream) BulkReadBytes(length int32) ([]byte, error) {
	return s.readBytes(length)
}

func (s *BitStream) ReadBits(mask uint32, bits int) uint32 {
	defer s.Advance(bits)
	return s.Peek(mask)
}

func (s *BitStream) readBytes(length int32) ([]byte, error) {
	if s.currentByteIndex+length > s.bufferSize {
		partialCount := s.currentByteIndex + length - s.bufferSize
		partial := s.byteBuffer[s.currentByteIndex:partialCount]

		s.advanceByteBufferIndex(partialCount)

		remainingCount := length - partialCount
		remaining := s.byteBuffer[s.currentByteIndex:remainingCount]

		return append(partial, remaining...), nil
	}

	result := s.byteBuffer[s.currentByteIndex : s.currentByteIndex+length]
	s.currentByteIndex += length

	return result, nil
}

func (s *BitStream) needsBufferRefresh(bytesToRead int32) bool {
	return s.currentByteIndex+bytesToRead > s.bufferSize
}

func (s *BitStream) refreshByteBuffer() (int, error) {
	bytesStillAvailable := s.bufferSize - s.currentByteIndex

	if bytesStillAvailable == s.bufferSize {
		bytesStillAvailable = 0
	}

	var intermediateBuffer []byte

	if bytesStillAvailable > 0 {
		intermediateBuffer = s.byteBuffer[s.currentByteIndex:]
	}

	newBuffer := make([]byte, s.bufferSize-bytesStillAvailable)
	bytesRead, err := io.ReadFull(*s.reader, newBuffer)

	if intermediateBuffer != nil {
		s.byteBuffer = append(intermediateBuffer, newBuffer...)
	} else {
		s.byteBuffer = newBuffer
	}

	s.currentByteIndex = 0

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		s.reachedEOF = true
		return bytesRead + int(bytesStillAvailable), nil
	}

	return bytesRead, err
}
