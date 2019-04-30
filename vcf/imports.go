package vcf

import (
	"io"
)

import (
	"github.com/brentp/vcfgo"
	"github.com/remeh/sizedwaitgroup"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

var DEBUG bool = false
var ONLYFIRST bool = false
var BREAKAT_THREAD int64 = 0
var BREAKAT_CHROM int64 = 0

//
// tools
var SliceIndex = tools.SliceIndex

//
// interfaces
type CallBackParameters = interfaces.CallBackParameters

//
// sized wait group
type SizedWaitGroup = sizedwaitgroup.SizedWaitGroup

//
// openfile
var OpenFile = openfile.OpenFile

//
// VCF
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
}

type VCFRegister = VCFRegisterRaw

type VCFCallBack func(*VCFSamples, *VCFRegister)
type VCFReaderType func(io.Reader, VCFCallBack, bool, []string)
type VCFMaskedReaderType = interfaces.VCFMaskedReaderType
type VCFMaskedReaderChromosomeType = interfaces.VCFMaskedReaderChromosomeType

// type VCFMaskedReaderType func(io.Reader, CallBackParameters)
// type VCFMaskedReaderChromosomeType func(io.Reader, bool, []string)
