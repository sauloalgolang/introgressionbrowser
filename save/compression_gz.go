package save

import (
	"compress/gzip"
	"io"
)

//
//
// GZip wrapper
//
//

//
// GZip Writer
//

// NewGzipWriter returns a new instance of a Gzip Writer
func NewGzipWriter(s io.Writer) GenericWriter {
	r, _ := gzip.NewWriterLevel(s, 1)

	n := &GzipWriterI{
		sn: r,
	}

	return n
}

// GzipWriterI holds a gzip writer
type GzipWriterI struct {
	sn *gzip.Writer
}

// Close closes the file handler
func (s *GzipWriterI) Close() error {
	return s.sn.Close()
}

// Flush flushes the file handler
func (s *GzipWriterI) Flush() error {
	return s.sn.Flush()
}

// Reset resets the file handler
func (s *GzipWriterI) Reset(w io.Writer) {
	s.sn.Reset(w)
}

// Write writes to file
func (s *GzipWriterI) Write(b []byte) (int, error) {
	return s.sn.Write(b)
}

//
// GZip Reader
//

// NewGzipReader returns a new instance of a Gzip Reader
func NewGzipReader(s io.Reader) GenericReader {
	r, _ := gzip.NewReader(s)

	n := &GzipReaderI{
		sn: r,
	}

	return n
}

// GzipReaderI holds a gzip reader
type GzipReaderI struct {
	sn *gzip.Reader
}

// Read reads from the file handler
func (s *GzipReaderI) Read(b []byte) (int, error) {
	return s.sn.Read(b)
}
