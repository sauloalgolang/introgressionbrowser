package save

import (
	"io"
)

//
// Generic Interfaces
//

type GenericWriter interface {
	Close() error
	Flush() error
	Reset(io.Writer)
	Write([]byte) (int, error)
}

type GenericReader interface {
	Read([]byte) (int, error)
}

type GenericNewWriter = func(io.Writer) GenericWriter
type GenericNewReader = func(io.Reader) GenericReader

type CompressorInterface struct {
	NewWriter GenericNewWriter
	NewReader GenericNewReader
}

//
// Implementations
//

var snappyInterface = CompressorInterface{
	NewWriter: NewSnappyWriter,
	NewReader: NewSnappyReader,
}

var gzipInterface = CompressorInterface{
	NewWriter: NewGzipWriter,
	NewReader: NewGzipReader,
}

var emptyInterface = CompressorInterface{
	NewWriter: nil,
	NewReader: nil,
}
