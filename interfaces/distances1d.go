package interfaces

import (
	"fmt"
	"math"
	"sync/atomic"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

//
//
// Matrix 1D
//
//
type DistanceMatrix1D struct {
	ChromosomeName string
	BlockSize      uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Dimension      uint64
	Size           uint64
	Data           DistanceRow
}

func NewDistanceMatrix1D(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64) *DistanceMatrix1D {
	size := dimension * (dimension - 1) / 2

	fmt.Println("   NewDistanceMatrix1D :: Chromosome: ", chromosomeName,
		" Dimension:", dimension,
		" Block Size: ", blockSize,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
		" Size:", size)

	r := DistanceMatrix1D{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		Dimension:      dimension,
		Size:           size,
		Data:           make(DistanceRow, size, size),
	}

	r.Clean()

	return &r
}

func (d *DistanceMatrix1D) add(e *DistanceMatrix1D, isAtomic bool) {
	if isAtomic {
		for i := range (*d).Data {
			atomic.AddUint64(&(*d).Data[i], atomic.LoadUint64(&(*e).Data[i]))
		}
	} else {
		for i := range (*d).Data {
			(*d).Data[i] += (*e).Data[i]
		}
	}
}

func (d *DistanceMatrix1D) Add(e *DistanceMatrix1D) {
	d.add(e, false)
}

func (d *DistanceMatrix1D) AddAtomic(e *DistanceMatrix1D) {
	d.add(e, true)
}

func (d *DistanceMatrix1D) Clean() {
	for i := range (*d).Data {
		(*d).Data[i] = uint64(0)
	}
}

// # https://stackoverflow.com/questions/27086195/linear-index-upper-triangular-matrix

func (d *DistanceMatrix1D) ijToK(i uint64, j uint64) uint64 {
	dim := float64(d.Dimension)
	fi := float64(i)
	fj := float64(j)

	fk := (dim * (dim - 1) / 2) - (dim-fi)*((dim-fi)-1)/2 + fj - fi - 1

	return uint64(fk)
}

func (d *DistanceMatrix1D) kToIJ(k uint64) (uint64, uint64) {
	dim := float64(d.Dimension)
	idx := float64(k)

	fi := dim - 2 - math.Floor(math.Sqrt(-8*idx+4*dim*(dim-1)-7)/2.0-0.5)
	fj := idx + fi + 1 - dim*(dim-1)/2 + (dim-fi)*((dim-fi)-1)/2

	return uint64(fi), uint64(fj)
}

func (d *DistanceMatrix1D) Set(p1 uint64, p2 uint64, val uint64) {
	(*d).Data[d.ijToK(p1, p2)] += val
}

func (d *DistanceMatrix1D) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	return (*d).Data[d.ijToK(p1, p2)]
}

func (d *DistanceMatrix1D) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix + "_matrix"
	fileName = save.GenFilename(baseName, format)
	return baseName, fileName
}

func (d *DistanceMatrix1D) Save(outPrefix string, format string) {
	baseName, _ := d.GenFilename(outPrefix, format)
	save.Save(baseName, format, d)
}
