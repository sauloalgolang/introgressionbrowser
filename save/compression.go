package save

import (
	"compress/gzip"
	"github.com/golang/snappy"
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
