package save

import (
	"github.com/golang/snappy"
	"io"
)

//
//
// Snappy wrapper
//
//

//
// Snappy Writer
//

// NewSnappyWriter generates a new instance of a snappy writer
func NewSnappyWriter(s io.Writer) GenericWriter {
	n := &SnappyWriterI{
		sn: snappy.NewBufferedWriter(s),
	}

	return n
}

// SnappyWriterI holds a snappy writer
type SnappyWriterI struct {
	sn *snappy.Writer
}

// Close closes file handler
func (s *SnappyWriterI) Close() error {
	return s.sn.Close()
}

// Flush flushes file handler
func (s *SnappyWriterI) Flush() error {
	return s.sn.Flush()
}

// Reset resets file handler
func (s *SnappyWriterI) Reset(w io.Writer) {
	s.sn.Reset(w)
}

// Write writes to file
func (s *SnappyWriterI) Write(b []byte) (int, error) {
	return s.sn.Write(b)
}

//
// Snappy Reader
//

// NewSnappyReader generates a new instance of a snappy reader
func NewSnappyReader(s io.Reader) GenericReader {
	n := &SnappyReaderI{
		sn: snappy.NewReader(s),
	}

	return n
}

// SnappyReaderI holds a snappy reader
type SnappyReaderI struct {
	sn *snappy.Reader
}

// Read reads from file handler
func (s *SnappyReaderI) Read(b []byte) (int, error) {
	return s.sn.Read(b)
}
