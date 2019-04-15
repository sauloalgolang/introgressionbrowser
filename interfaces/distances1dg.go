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

// type DistanceRow32 []uint32
// type DistanceRow64 []uint64

type DistanceMatrix1Dg struct {
	ChromosomeName string
	BlockSize      uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Dimension      uint64
	Size           uint64
	Bits           int
	Data32         DistanceRow32
	Data64         DistanceRow64
	// Data           []interface{}
}

func NewDistanceMatrix1Dg32(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, blockPosition, blockNumber, dimension, 32)
}

func NewDistanceMatrix1Dg64(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, blockPosition, blockNumber, dimension, 64)
}

func NewDistanceMatrix1Dg(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, dimension uint64, bits int) *DistanceMatrix1Dg {
	size := dimension * (dimension - 1) / 2

	fmt.Println("   NewDistanceMatrix1D :: Chromosome: ", chromosomeName,
		" Dimension:", dimension,
		" Block Size: ", blockSize,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
		" Size:", size,
		" Bits:", bits,
	)

	r := DistanceMatrix1Dg{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		Dimension:      dimension,
		Size:           size,
		Bits:           bits,
	}

	if r.Bits == 32 {
		// r.Data = make(DistanceRow32, size, size)
		r.Data32 = make(DistanceRow32, size, size)
		r.Data64 = make(DistanceRow64, 0, 0)
	} else if r.Bits == 64 {
		// r.Data = make(DistanceRow64, size, size)
		r.Data32 = make(DistanceRow32, 0, 0)
		r.Data64 = make(DistanceRow64, size, size)
	}

	r.Clean()

	return &r
}

//
// Exported Methods
//

func (d *DistanceMatrix1Dg) Add(e *DistanceMatrix1Dg) {
	d.add(e, false)
}

func (d *DistanceMatrix1Dg) AddAtomic(e *DistanceMatrix1Dg) {
	d.add(e, true)
}

func (d *DistanceMatrix1Dg) Clean() {
	if d.Bits == 32 {
		d.clean32()
	} else if d.Bits == 64 {
		d.clean64()
	}
}

func (d *DistanceMatrix1Dg) clean32() {
	for i := range (*d).Data32 {
		(*d).Data32[i] = uint32(0)
	}
}

func (d *DistanceMatrix1Dg) clean64() {
	for i := range (*d).Data64 {
		(*d).Data64[i] = uint64(0)
	}
}

func (d *DistanceMatrix1Dg) Set(p1 uint64, p2 uint64, val uint64) {
	p := d.ijToK(p1, p2)

	if d.Bits == 32 {
		d.set32(p, val)
	} else if d.Bits == 64 {
		d.set64(p, val)
	}
}

func (d *DistanceMatrix1Dg) set32(p uint64, val uint64) {
	v := (*d).Data32[p]
	r := v + uint32(val)

	if val >= uint64(math.MaxUint32) {
		fmt.Println("count overflow")
		os.Exit(1)
	}

	(*d).Data32[p] = r
}

func (d *DistanceMatrix1Dg) set64(p uint64, val uint64) {
	(*d).Data64[p] = val
}

func (d *DistanceMatrix1Dg) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	p := d.ijToK(p1, p2)

	if d.Bits == 32 {
		return uint64((*d).Data32[p])
	} else if d.Bits == 64 {
		return (*d).Data64[p]
	}

	return 0
}

func (d *DistanceMatrix1Dg) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "_matrix"

	saver := save.NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Unexported Methods
//

func (d *DistanceMatrix1Dg) ijToK(i uint64, j uint64) uint64 {
	return ijToK(d.Dimension, i, j)
}

func (d *DistanceMatrix1Dg) kToIJ(k uint64) (uint64, uint64) {
	return kToIJ(d.Dimension, k)
}

func (d *DistanceMatrix1Dg) add(e *DistanceMatrix1Dg, isAtomic bool) {
	if d.Bits == 32 {
		d.add32(e, isAtomic)
	} else if d.Bits == 64 {
		d.add64(e, isAtomic)
	}
}

func (d *DistanceMatrix1Dg) add32(e *DistanceMatrix1Dg, isAtomic bool) {
	if isAtomic {
		for i := range (*d).Data32 {
			atomic.AddUint32(&(*d).Data32[i], atomic.LoadUint32(&(*e).Data32[i]))
		}
	} else {
		mi := uint64(math.MaxInt32)
		for i := range (*d).Data32 {
			if uint64((*d).Data32[i])+uint64((*e).Data32[i]) >= mi {
				fmt.Println("counter overflow")
				os.Exit(1)
			}
			(*d).Data32[i] += (*e).Data32[i]
		}
	}

}

func (d *DistanceMatrix1Dg) add64(e *DistanceMatrix1Dg, isAtomic bool) {
	if isAtomic {
		for i := range (*d).Data64 {
			atomic.AddUint64(&(*d).Data64[i], atomic.LoadUint64(&(*e).Data64[i]))
		}
	} else {
		for i := range (*d).Data64 {
			(*d).Data64[i] += (*e).Data64[i]
		}
	}

}

//
// Save and Load
//
func (d *DistanceMatrix1Dg) Save(outPrefix string, format string, compression string) {
	d.saveLoad(true, outPrefix, format, compression)
}

func (d *DistanceMatrix1Dg) Load(outPrefix string, format string, compression string) {
	d.saveLoad(false, outPrefix, format, compression)
}

func (d *DistanceMatrix1Dg) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := d.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving block matrix    : ", outPrefix, " block num: ", d.BlockNumber)
		saver.Save(d)
	} else {
		fmt.Println("loading block matrix   : ", outPrefix, " block num: ", d.BlockNumber)
		saver.Load(d)
	}
}
