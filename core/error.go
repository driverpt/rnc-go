package core

import "fmt"

type RNCError interface {
	error
}

type NonRNCStream struct {
	error RNCError
}

func (e NonRNCStream) Error() string {
	return fmt.Sprint("Stream is not RNC Encoded")
}

type UnsupportedCompressionMethod struct {
	error   RNCError
	version uint8
}

func (e UnsupportedCompressionMethod) Error() string {
	return fmt.Sprintf("Compression Method: RNC%d not supported", e.version)
}

func NewUnsupportedCompressionMethodError(version uint8) *UnsupportedCompressionMethod {
	return &UnsupportedCompressionMethod{
		version: version,
	}
}

type FileSizeMismatch struct {
	expected int32
	current  int32
}

func (e FileSizeMismatch) Error() string {
	return fmt.Sprintf("Invalid File Size - expected %d, current %d", e.expected, e.current)
}

func NewFileSizeMismatchError(current int32, expected int32) *FileSizeMismatch {
	return &FileSizeMismatch{
		expected: expected,
		current:  current,
	}
}

type UnpackedCrcError struct {
	expected uint16
	current  uint16
}

func NewUnpackedCrcError(current uint16, expected uint16) *UnpackedCrcError {
	return &UnpackedCrcError{
		expected: expected,
		current:  current,
	}
}

func (e UnpackedCrcError) Error() string {
	return fmt.Sprintf("Invalid Unpacked CRC Error - expected %d, current %d", e.expected, e.current)
}

type PackedCrcError struct {
	expected uint16
	current  uint16
}

func NewPackedCrcError(current uint16, expected uint16) *PackedCrcError {
	return &PackedCrcError{
		expected: expected,
		current:  current,
	}
}

func (e *PackedCrcError) Error() string {
	return fmt.Sprintf("Invalid Packed CRC Error - expected %d, current %d", e.expected, e.current)
}
