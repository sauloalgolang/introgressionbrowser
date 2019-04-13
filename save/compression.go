package save

import (
	"compress/gzip"
	"github.com/golang/snappy"
)

type CompressorInterface struct {
	NewWriter interface{}
	NewReader interface{}
}

var emptyInterface = CompressorInterface{
	NewWriter: nil,
	NewReader: nil,
}

var snappyInterface = CompressorInterface{
	NewWriter: snappy.NewWriter,
	NewReader: snappy.NewReader,
}

var gzipInterface = CompressorInterface{
	NewWriter: gzip.NewWriter,
	NewReader: gzip.NewReader,
}
