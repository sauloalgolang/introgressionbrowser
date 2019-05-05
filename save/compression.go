package save

import (
	"io"
)

//
// Generic Interfaces
//

// GenericWriter interface for writers
type GenericWriter interface {
	Close() error
	Flush() error
	Reset(io.Writer)
	Write([]byte) (int, error)
}

// GenericReader interface for readers
type GenericReader interface {
	Read([]byte) (int, error)
}

// GenericNewWriter alias for a generic writer
type GenericNewWriter = func(io.Writer) GenericWriter

// GenericNewReader alias for a generic reader
type GenericNewReader = func(io.Reader) GenericReader

// CompressorInterface holds a writer and a reader for compressor
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
