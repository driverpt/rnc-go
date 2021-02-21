package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/driverpt/rnc-go/factory"
	"io"
	"log"
	"os"
)

func main() {
	crcOnly := flag.Bool("crc-only", false, "Only check for Packed CRC")
	flag.Parse()

	if flag.NArg() != 1 {
		exitAndPrintUsage(nil)
	}

	filePath := flag.Args()[0]

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		exitAndPrintUsage(&err)
	}

	header, err := factory.ReadRNCHeader(file)
	if err != nil {
		exitAndPrintUsage(&err)
	}

	// TODO: Probably we need to have a separate Checksum Verifier
	seek(file, 0)
	_, err = factory.VerifyPackedChecksum(header, file)
	if err != nil {
		exitAndPrintUsage(&err)
	}

	if *crcOnly {
		os.Exit(0)
	}

	seek(file, 0)
	reader, err := factory.NewRNCReader(file)

	if err != nil {
		exitAndPrintUsage(&err)
	}

	result, err := reader.Unpack()

	if err != nil {
		exitAndPrintUsage(&err)
	}

	byteStream := bytes.NewReader(result)

	_, err = io.Copy(os.Stdout, byteStream)
	factory.VerifyUnpackedChecksum(*header, byteStream)

	if err != nil {
		exitAndPrintUsage(&err)
	}

	os.Exit(0)
}

func seek(file *os.File, position int32) {
	_, err := file.Seek(int64(position), io.SeekStart)

	if err != nil {
		log.Fatal(err)
	}
}

func exitAndPrintUsage(error *error) {
	name, _ := os.Executable()

	fmt.Printf("Usage: %s [file_to_unpack]", name)
	flag.PrintDefaults()
	if error != nil {
		log.Fatal(error)
	}
}
