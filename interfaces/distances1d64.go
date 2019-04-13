package interfaces

import (
	"fmt"
	"sync/atomic"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

//
//
// Matrix 1D
//
//
type DistanceMatrix1D64 struct {
	ChromosomeName string
	BlockSize      uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Dimension      uint64
	Size           uint64
	Data           DistanceRow
}

func NewDistanceMatrix1D64(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64) *DistanceMatrix1D64 {
	size := dimension * (dimension - 1) / 2

	fmt.Println("   NewDistanceMatrix1D64 :: Chromosome: ", chromosomeName,
		" Dimension:", dimension,
		" Block Size: ", blockSize,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
		" Size:", size)

	r := DistanceMatrix1D64{
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

func (d *DistanceMatrix1D64) add(e *DistanceMatrix1D64, isAtomic bool) {
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

func (d *DistanceMatrix1D64) Add(e *DistanceMatrix1D64) {
	d.add(e, false)
}

func (d *DistanceMatrix1D64) AddAtomic(e *DistanceMatrix1D64) {
	d.add(e, true)
}

func (d *DistanceMatrix1D64) Clean() {
	for i := range (*d).Data {
		(*d).Data[i] = uint64(0)
	}
}

func (d *DistanceMatrix1D64) ijToK(i uint64, j uint64) uint64 {
	return ijToK(d.Dimension, i, j)
}

func (d *DistanceMatrix1D64) kToIJ(k uint64) (uint64, uint64) {
	return kToIJ(d.Dimension, k)
}

func (d *DistanceMatrix1D64) Set(p1 uint64, p2 uint64, val uint64) {
	(*d).Data[d.ijToK(p1, p2)] += val
}

func (d *DistanceMatrix1D64) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	return (*d).Data[d.ijToK(p1, p2)]
}

func (d *DistanceMatrix1D64) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix + "_matrix"

	saver := save.NewSaver(baseName, format)

	fileName = saver.GenFilename()

	return baseName, fileName
}

func (d *DistanceMatrix1D64) Save(outPrefix string, format string) {
	baseName, _ := d.GenFilename(outPrefix, format)

	saver := save.NewSaver(baseName, format)

	saver.Save(d)
}

func (d *DistanceMatrix1D64) Load(outPrefix string, format string) {
	baseName, _ := d.GenFilename(outPrefix, format)

	saver := save.NewSaver(baseName, format)
	saver.Load(d)
}
