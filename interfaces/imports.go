package interfaces

import (
	"io"
)

import (
	"github.com/brentp/vcfgo"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/matrix"
	"github.com/sauloalgolang/introgressionbrowser/save"
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
	LineNumber       int64
	Chromosome       string
	ChromosomeNumber int
	Position         uint64
	Alt              []string
	Samples          VCFSamplesGT
	Distance         *DistanceMatrix
	TempDistance     *DistanceMatrix
	// SampleNames []string
	// IsHomozygous bool
	// IsIndel      bool
	// IsMNP        bool
	// Quality      float32
	// Info         map[string]interface{}
	// Filter       string
	// NumAlt       uint64
	// Phased       bool
	// Fields       map[string]string
}

type VCFRegister = VCFRegisterRaw

type VCFCallBack func(*VCFSamples, *VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool, []string)
type VCFMaskedReaderType func(io.Reader, CallBackParameters)
type VCFMaskedReaderChromosomeType func(io.Reader, bool, []string)

//
// Distance
//

type DistanceMatrix = matrix.DistanceMatrix
type DistanceTable = matrix.DistanceTable

var NewSaver = save.NewSaver
