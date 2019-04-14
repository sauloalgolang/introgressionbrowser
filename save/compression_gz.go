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

func NewGzipWriter(s io.Writer) GenericWriter {
	r, _ := gzip.NewWriterLevel(s, 1)

	n := &GzipWriterI{
		sn: r,
	}

	return n
}

type GzipWriterI struct {
	sn *gzip.Writer
}

func (s *GzipWriterI) Close() error {
	return s.sn.Close()
}

func (s *GzipWriterI) Flush() error {
	return s.sn.Flush()
}

func (s *GzipWriterI) Reset(w io.Writer) {
	s.sn.Reset(w)
}

func (s *GzipWriterI) Write(b []byte) (int, error) {
	return s.sn.Write(b)
}

//
// GZip Reader
//

func NewGzipReader(s io.Reader) GenericReader {
	r, _ := gzip.NewReader(s)

	n := &GzipReaderI{
		sn: r,
	}

	return n
}

type GzipReaderI struct {
	sn *gzip.Reader
}

func (s *GzipReaderI) Read(b []byte) (int, error) {
	return s.sn.Read(b)
}
