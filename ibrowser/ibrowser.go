package IBrowser

import (
	"fmt"
	"io"
)

type VCFCallBack func(*VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool)
type VCFMaskedReaderType func(io.Reader, bool)

type VCFRegister struct {
	someint int
}

type IBrowser struct {
	reader VCFReaderType
}

func NewIBrowser(reader VCFReaderType) *IBrowser {
	ib := new(IBrowser)
	ib.reader = reader
	return ib
}

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool) {
	ib.reader(r, ib.RegisterCallBack, continueOnError)
}

func (ib *IBrowser) RegisterCallBack(reg *VCFRegister) {
	fmt.Println("got register", reg)
}
