package rnc1

import (
	"errors"
	"github.com/driverpt/rnc-go/core"
	"io"
)


type RNC1Reader struct {
	Header *core.RNCHeader
	Reader *io.Reader
	initialOffset int32
}

func (r RNC1Reader) Unpack() ([]byte, error) {
	// Assume that Header has already been read
	r.initialOffset = core.HeaderSizeInBytes
	outputOffset := int32(0)
	result := make([]byte, r.Header.OriginalSize)

	bitstream := NewBitStream(*r.Reader, 0, int32(r.Header.CompressedSize))
	// First 2 bits are ignored
	bitstream.Advance(2)

	for outputOffset < int32(r.Header.OriginalSize) {
		literalDataSizes := ReadHuffmanTable(&bitstream)
		distance := ReadHuffmanTable(&bitstream)
		length := ReadHuffmanTable(&bitstream)

		subChunks := bitstream.ReadBits(0xFFFF, 16)

		for true {
			huffmanLength := ReadHuffman(&literalDataSizes, &bitstream)
			if huffmanLength == -1 {
				return nil, errors.New("Huffman Decode Error")
			}

			if huffmanLength != 0 {
				buffer, err := bitstream.BulkReadBytes(huffmanLength)
				if err != nil {
					return nil, err
				}

				for i := int32(0); i < huffmanLength; i++ {
					result[i + outputOffset] = buffer[i]
				}

				outputOffset += huffmanLength
				bitstream.RefreshBuffer()
			}

			subChunks -= 1
			if subChunks <= 0 {
				break
			}

			pos := ReadHuffman(&distance, &bitstream)
			if pos == -1 {
				panic("Huffman Decode Error")
			}

			pos += 1

			huffLength := ReadHuffman(&length, &bitstream)

			if huffLength == -1 {
				panic("Huffman Decode Error")
			}

			huffLength += 2

			for huffLength > 0 {
				result[outputOffset] = result[outputOffset-pos]
				huffLength--
				outputOffset++
			}
		}
	}

	if outputOffset != int32(r.Header.OriginalSize) {
		return nil, core.NewFileSizeMismatchError(outputOffset, int32(r.Header.OriginalSize))
	}

	return result, nil
}
