## Rob Northen Compression Decoder in Go

This implementation is inspired from CorsixTH project `rnc.cpp`

### Usage
```go
import "github.com/driverpt/rnc-go/factory"

func main() {
    file, err := os.Open(**RNC Compressed File**)
    
    // If you want checksum verification
    header, err := factory.ReadRNCHeader(file)
    _, err = factory.VerifyPackedChecksum(header, file)

    // Ensure that you put the file back into 0 if you do Checksum
    _, err := file.Seek(0, io.SeekStart)

    reader, err := factory.NewRNCReader(file)
    
    // Resulting file is stored in here
    var result byte[] = reader.Unpack()
    byteStream := bytes.NewReader(result)
    _, err := factory.VerifyUnpackedChecksum(header, byteStream)
}

```


#### TODO
 - Support RNC2
 - Add Tests
