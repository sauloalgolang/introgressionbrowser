package interfaces

import (
	"io"
)

import (
	"github.com/brentp/vcfgo"
)

//
//
// VCF Register
//
//

type VCFRegisterVcfGo = vcfgo.Variant
type VCFSamples = []string
type VCFGTVal []int
type VCFGT struct {
	GT VCFGTVal
}
type VCFSamplesGT = []VCFGT

type VCFRegisterRaw struct {
	// SampleNames []string
	// IsHomozygous bool
	// IsIndel      bool
	// IsMNP        bool
	LineNumber int64
	Chromosome string
	Position   uint64
	// Quality      float32
	// Info         map[string]interface{}
	// Filter       string
	// NumAlt       uint64
	// Phased       bool
	Alt     []string
	Samples VCFSamplesGT
	// Fields       map[string]string
	Distance     *DistanceMatrix
	TempDistance *DistanceMatrix
}

type VCFRegister = VCFRegisterRaw

type VCFCallBack func(*VCFSamples, *VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool, []string)
type VCFMaskedReaderType func(io.Reader, bool)
type VCFMaskedReaderChromosomeType func(io.Reader, bool, []string)
