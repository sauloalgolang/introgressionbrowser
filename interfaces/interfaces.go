package interfaces

import (
	"io"
)

import (
	"github.com/brentp/vcfgo"
)

type VCFRegister = vcfgo.Variant
type VCFSamples = []string

type VCFCallBack func(*VCFSamples, *VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool)
type VCFMaskedReaderType func(io.Reader, bool)

// type VCFRegister struct {
// 	Samples      *[]string
// 	IsHomozygous bool
// 	IsIndel      bool
// 	IsMNP        bool
// 	Row          uint64
// 	Chromosome   string
// 	Position     uint64
// 	Quality      float32
// 	Info         map[string]interface{}
// 	Filter       string
// 	NumAlt       uint64
// 	Phased       bool
// 	GT           [][]int
// 	Fields       map[string]string
// }
