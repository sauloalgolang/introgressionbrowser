package interfaces

import (
	"io"
)

type VCFCallBack func(*VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool)
type VCFMaskedReaderType func(io.Reader, bool)

type VCFRegister struct {
	someint int
}
