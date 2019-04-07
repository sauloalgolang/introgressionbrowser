package interfaces

import (
	"io"
	"sync/atomic"
)

import (
	"github.com/brentp/vcfgo"
)

type DistanceRow []uint64
type DistanceMatrix [][]uint64
type DistanceTable []uint64

func (d *DistanceMatrix) add(e *DistanceMatrix, isAtomic bool) {
	for i := range *d {
		di := &(*d)[i]
		ei := &(*e)[i]

		for j := i + 1; j < len(*d); j++ {
			if isAtomic {
				atomic.AddUint64(&(*di)[j], atomic.LoadUint64(&(*ei)[j]))

			} else {
				(*di)[j] += (*ei)[j]
			}
		}
	}
}

func (d *DistanceMatrix) Add(e *DistanceMatrix) {
	d.add(e, false)
}
func (d *DistanceMatrix) AddAtomic(e *DistanceMatrix) {
	d.add(e, true)
}

func (d *DistanceMatrix) Clean() {
	for i := range *d {
		ti := &(*d)[i]
		for j := i + 1; j < len(*d); j++ {
			(*ti)[j] = uint64(0)
		}
	}
}

func (d *DistanceMatrix) Set(p1 uint64, p2 uint64, val uint64) {
	(*d)[p1][p2] += val
}

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
type VCFReaderType func(io.Reader, VCFCallBack, bool, string)
type VCFMaskedReaderType func(io.Reader, bool)
type VCFMaskedReaderChromosomeType func(io.Reader, bool, string)
