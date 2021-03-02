package factory

import (
	"github.com/driverpt/rnc-go/core"
	"github.com/driverpt/rnc-go/rnc1"
	"io"
)

func ReadRNCHeader(src io.Reader) (*core.RNCHeader, error) {
	return core.ReadHeader(src)
}

func NewRNCReader(src io.Reader) (core.RNCReader, error) {
	header, err := core.ReadHeader(src)

	if err != nil {
		return nil, err
	}

	switch header.CompressionMethod {
	case 1:
		return rnc1.RNC1Reader{
			Header: header,
			Reader: &src,
		}, nil
	case 2:
		// Add RNC2 here
	}

	return nil, core.NewUnsupportedCompressionMethodError(header.CompressionMethod)
}

func VerifyPackedChecksum(header *core.RNCHeader, src io.Reader) (bool, error) {
	switch header.CompressionMethod {
	case 1:
		checksum, err := rnc1.ComputeChecksum(src, core.HeaderSizeInBytes, int32(header.CompressedSize)+core.HeaderSizeInBytes)
		if err != nil {
			return false, err
		}

		if checksum != header.CompressedCRC {
			return false, core.NewPackedCrcError(checksum, header.CompressedCRC)
		}
		return true, nil
	case 2:
		// Add RNC2 here
	}

	return false, core.NewUnsupportedCompressionMethodError(header.CompressionMethod)
}

func VerifyUnpackedChecksum(header core.RNCHeader, src io.Reader) (bool, error) {
	switch header.CompressionMethod {
	case 1:
		checksum, err := rnc1.ComputeChecksum(src, 0, int32(header.OriginalSize))
		if err != nil {
			return false, err
		}

		if checksum != header.OriginalCRC {
			return false, core.NewUnpackedCrcError(checksum, header.OriginalCRC)
		}
		return true, nil
	case 2:
		// Add RNC2 here
	}

	return false, core.NewUnsupportedCompressionMethodError(header.CompressionMethod)
}
