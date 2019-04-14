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

func NewSnappyWriter(s io.Writer) GenericWriter {
	n := &SnappyWriterI{
		sn: snappy.NewBufferedWriter(s),
	}

	return n
}

type SnappyWriterI struct {
	sn *snappy.Writer
}

func (s *SnappyWriterI) Close() error {
	return s.sn.Close()
}

func (s *SnappyWriterI) Flush() error {
	return s.sn.Flush()
}

func (s *SnappyWriterI) Reset(w io.Writer) {
	s.sn.Reset(w)
}

func (s *SnappyWriterI) Write(b []byte) (int, error) {
	return s.sn.Write(b)
}

//
// Snappy Reader
//

func NewSnappyReader(s io.Reader) GenericReader {
	n := &SnappyReaderI{
		sn: snappy.NewReader(s),
	}

	return n
}

type SnappyReaderI struct {
	sn *snappy.Reader
}

func (s *SnappyReaderI) Read(b []byte) (int, error) {
	return s.sn.Read(b)
}
