package ibrowser

import (
	"fmt"
	"io"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

type IBrowser struct {
	reader interfaces.VCFReaderType
}

func NewIBrowser(reader interfaces.VCFReaderType) *IBrowser {
	ib := new(IBrowser)
	ib.reader = reader
	return ib
}

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool) {
	ib.reader(r, ib.RegisterCallBack, continueOnError)
}

func (ib *IBrowser) RegisterCallBack(reg *interfaces.VCFRegister) {
	fmt.Println("got register", reg)
}
