package vcf

import (
	"io"
)

import (
	// "github.com/brentp/vcfgo"
	"github.com/remeh/sizedwaitgroup"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

// Debug defines whether to debug or not
var Debug = false

// OnlyFirst defines to only read the first scaffold
var OnlyFirst = false

// BreakAtThread Number of row to read per thread
var BreakAtThread = int64(0)

// BreakAtChrom Number of rows to read per chroms
var BreakAtChrom = int64(0)

// SliceIndex alias to tools.SliceIndex
var SliceIndex = tools.SliceIndex

// CallBackParameters alias to interfaces.CallBackParameters
type CallBackParameters = interfaces.CallBackParameters

// SizedWaitGroup alias to sizedwaitgroup.SizedWaitGroup
type SizedWaitGroup = sizedwaitgroup.SizedWaitGroup

// OpenAnyFile alias to openfile.OpenFile
var OpenAnyFile = openfile.OpenFile

//
// VCF
//

// Samples alias to []string - list of sample names
type Samples = []string

// GenotypeVal values from genotype call
type GenotypeVal []int

// RegisterGenotype struc holding the values from a genotype call
type RegisterGenotype struct {
	Genotype GenotypeVal
}

// SamplesGT list of genotypes for each sample
type SamplesGenotype = []RegisterGenotype

// RegisterRaw holds a vcf register
type RegisterRaw struct {
	LineNumber       int64
	Chromosome       string
	ChromosomeNumber int
	Position         uint64
	Alt              []string
	Samples          SamplesGenotype
	Distance         *DistanceMatrix
	TempDistance     *DistanceMatrix
}

// Register alias to the default register
type Register = RegisterRaw

// MaskedReaderType alias to interfaces.MaskedReaderType
type MaskedReaderType = interfaces.VCFMaskedReaderType

// MaskedReaderChromosomeType alias to interfaces.MaskedReaderChromosomeType
type MaskedReaderChromosomeType = interfaces.VCFMaskedReaderChromosomeType

// RegisterCallBack callback when finding a register
type RegisterCallBack func(*Samples, *Register)

// ReaderType type for a file reader
type ReaderType func(io.Reader, RegisterCallBack, bool, []string)

// type MaskedReaderType func(io.Reader, CallBackParameters)
// type MaskedReaderChromosomeType func(io.Reader, bool, []string)
