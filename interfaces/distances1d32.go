package interfaces

import (
	"fmt"
	"math"
	"os"
	"sync/atomic"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

//
//
// Matrix 1D
//
//
type DistanceRow32 []uint32

type DistanceMatrix1D32 struct {
	ChromosomeName string
	BlockSize      uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Dimension      uint64
	Size           uint64
	Data           DistanceRow32
}

func NewDistanceMatrix1D32(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64) *DistanceMatrix1D32 {
	size := dimension * (dimension - 1) / 2

	fmt.Println("   NewDistanceMatrix1D32 :: Chromosome: ", chromosomeName,
		" Dimension:", dimension,
		" Block Size: ", blockSize,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
		" Size:", size)

	r := DistanceMatrix1D32{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		Dimension:      dimension,
		Size:           size,
		Data:           make(DistanceRow32, size, size),
	}

	r.Clean()

	return &r
}

func (d *DistanceMatrix1D32) add(e *DistanceMatrix1D32, isAtomic bool) {
	if isAtomic {
		for i := range (*d).Data {
			atomic.AddUint32(&(*d).Data[i], atomic.LoadUint32(&(*e).Data[i]))
		}
	} else {
		for i := range (*d).Data {
			(*d).Data[i] += (*e).Data[i]
		}
	}
}

func (d *DistanceMatrix1D32) Add(e *DistanceMatrix1D32) {
	d.add(e, false)
}

func (d *DistanceMatrix1D32) AddAtomic(e *DistanceMatrix1D32) {
	d.add(e, true)
}

func (d *DistanceMatrix1D32) Clean() {
	for i := range (*d).Data {
		(*d).Data[i] = uint32(0)
	}
}

func (d *DistanceMatrix1D32) ijToK(i uint64, j uint64) uint64 {
	return ijToK(d.Dimension, i, j)
}

func (d *DistanceMatrix1D32) kToIJ(k uint64) (uint64, uint64) {
	return kToIJ(d.Dimension, k)
}

func (d *DistanceMatrix1D32) Set(p1 uint64, p2 uint64, val uint64) {
	p := d.ijToK(p1, p2)
	v := (*d).Data[p]
	r := v + uint32(val)

	if uint64(v)+val >= uint64(math.MaxUint32) {
		fmt.Println("count overflow")
		os.Exit(1)
	}

	(*d).Data[p] = r
}

func (d *DistanceMatrix1D32) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	return uint64((*d).Data[d.ijToK(p1, p2)])
}

func (d *DistanceMatrix1D32) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix + "_matrix"

	saver := save.NewSaver(baseName, format)

	fileName = saver.GenFilename()

	return baseName, fileName
}

func (d *DistanceMatrix1D32) Save(outPrefix string, format string) {
	baseName, _ := d.GenFilename(outPrefix, format)

	saver := save.NewSaver(baseName, format)
	saver.Save(d)
}

func (d *DistanceMatrix1D32) Load(outPrefix string, format string) {
	baseName, _ := d.GenFilename(outPrefix, format)

	saver := save.NewSaver(baseName, format)
	saver.Load(d)
}
