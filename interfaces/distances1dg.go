package interfaces

import (
	"fmt"
	"math"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

//
//
// Matrix 1D
//
//

type DistanceRow16 []uint16
type DistanceRow32 []uint32
type DistanceRow64 []uint64

type DistanceMatrix1Dg struct {
	ChromosomeName string
	BlockSize      uint64
	Dimension      uint64
	Size           uint64
	BlockPosition  uint64
	BlockNumber    uint64
	NumBits        int
	Data16         DistanceRow16
	Data32         DistanceRow32
	Data64         DistanceRow64
	// Data           []interface{}
}

func NewDistanceMatrix1Dg16(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 16, dimension, blockPosition, blockNumber)
}

func NewDistanceMatrix1Dg32(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 32, dimension, blockPosition, blockNumber)
}

func NewDistanceMatrix1Dg64(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 64, dimension, blockPosition, blockNumber)
}

func NewDistanceMatrix1Dg(chromosomeName string, blockSize uint64, numBits int, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	size := dimension * (dimension - 1) / 2

	fmt.Println("    NewDistanceMatrix1D :: Chromosome: ", chromosomeName,
		" Block Size: ", blockSize,
		" Bits:", numBits,
		" Dimension:", dimension,
		" Size:", size,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
	)

	d := DistanceMatrix1Dg{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		Dimension:      dimension,
		Size:           size,
		NumBits:        numBits,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
	}

	if d.NumBits == 16 {
		// d.Data = make(DistanceRow32, size, size)
		d.Data16 = make(DistanceRow16, size, size)
		d.Data32 = make(DistanceRow32, 0, 0)
		d.Data64 = make(DistanceRow64, 0, 0)
	} else if d.NumBits == 32 {
		// d.Data = make(DistanceRow32, size, size)
		d.Data16 = make(DistanceRow16, 0, 0)
		d.Data32 = make(DistanceRow32, size, size)
		d.Data64 = make(DistanceRow64, 0, 0)
	} else if d.NumBits == 64 {
		// d.Data = make(DistanceRow64, size, size)
		d.Data16 = make(DistanceRow16, 0, 0)
		d.Data32 = make(DistanceRow32, 0, 0)
		d.Data64 = make(DistanceRow64, size, size)
	}

	d.Clean()

	return &d
}

//
// Clean
//

func (d *DistanceMatrix1Dg) Clean() {
	if d.NumBits == 16 {
		d.clean16()
	} else if d.NumBits == 32 {
		d.clean32()
	} else if d.NumBits == 64 {
		d.clean64()
	}
}

func (d *DistanceMatrix1Dg) clean16() {
	for i := range (*d).Data16 {
		(*d).Data16[i] = uint16(0)
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

//
// Set
//

func (d *DistanceMatrix1Dg) Set(p1 uint64, p2 uint64, val uint64) {
	p := d.ijToK(p1, p2)

	if d.NumBits == 16 {
		d.set16(p, val)
	} else if d.NumBits == 32 {
		d.set32(p, val)
	} else if d.NumBits == 64 {
		d.set64(p, val)
	}
}

func (d *DistanceMatrix1Dg) set16(p uint64, val uint64) {
	v := (*d).Data16[p]
	r := v + uint16(val)

	if val >= uint64(math.MaxUint16) {
		fmt.Println("count overflow")
		os.Exit(1)
	}

	(*d).Data16[p] = r
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

//
// Add
//

func (d *DistanceMatrix1Dg) Add(e *DistanceMatrix1Dg) {
	d.add(e)
}

func (d *DistanceMatrix1Dg) add(e *DistanceMatrix1Dg) {
	if d.NumBits == 16 {
		d.add16(e)
	} else if d.NumBits == 32 {
		d.add32(e)
	} else if d.NumBits == 64 {
		d.add64(e)
	}
}

func (d *DistanceMatrix1Dg) add16(e *DistanceMatrix1Dg) {
	mi := uint64(math.MaxInt16)
	for i := range (*d).Data16 {
		if uint64((*d).Data16[i])+uint64((*e).Data16[i]) >= mi {
			fmt.Println("counter overflow")
			os.Exit(1)
		}
		(*d).Data16[i] += (*e).Data16[i]
	}
}

func (d *DistanceMatrix1Dg) add32(e *DistanceMatrix1Dg) {
	mi := uint64(math.MaxInt32)
	for i := range (*d).Data32 {
		if uint64((*d).Data32[i])+uint64((*e).Data32[i]) >= mi {
			fmt.Println("counter overflow")
			os.Exit(1)
		}
		(*d).Data32[i] += (*e).Data32[i]
	}
}

func (d *DistanceMatrix1Dg) add64(e *DistanceMatrix1Dg) {
	for i := range (*d).Data64 {
		(*d).Data64[i] += (*e).Data64[i]
	}
}

//
// IsEqual
//

func (d *DistanceMatrix1Dg) IsEqual(e *DistanceMatrix1Dg) (res bool) {
	res = true

	// res = res && (d.ChromosomeName == e.ChromosomeName)
	//
	// if !res {
	// 	fmt.Printf("IsEqual :: Failed matrix %s - #%d check - ChromosomeName %s != %s\n", d.ChromosomeName, d.BlockNumber, d.ChromosomeName, e.ChromosomeName)
	// 	return res
	// }

	res = res && (d.BlockSize == e.BlockSize)

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check - BlockSize %d != %d\n", d.ChromosomeName, d.BlockNumber, d.BlockSize, e.BlockSize)
		return res
	}

	res = res && (d.Dimension == e.Dimension)

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check - Dimension %d != %d\n", d.ChromosomeName, d.BlockNumber, d.Dimension, e.Dimension)
		return res
	}

	res = res && (d.NumBits == e.NumBits)

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check - NumBits %d != %d\n", d.ChromosomeName, d.BlockNumber, d.NumBits, e.NumBits)
		return res
	}

	res = res && (d.Size == e.Size)

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check - Size %d != %d\n", d.ChromosomeName, d.BlockNumber, d.Size, e.Size)
		return res
	}

	if d.NumBits == 16 {
		d.check16(e)
	} else if d.NumBits == 32 {
		d.check32(e)
	} else if d.NumBits == 64 {
		d.check64(e)
	}

	return res

}

func (d *DistanceMatrix1Dg) check16(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.Data16)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 16 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.Data16)))
		return res
	}

	res = res && (e.Size == uint64(len(e.Data16)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 16 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.Data16)))
		return res
	}

	for i := range (*d).Data16 {
		res = res && ((*d).Data16[i] == (*e).Data16[i])

		if !res {
			fmt.Printf("IsEqual :: Failed matrix %s - #%d check 16 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).Data16[i], (*e).Data16[i])
		}
	}

	return res
}

func (d *DistanceMatrix1Dg) check32(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.Data32)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 32 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.Data32)))
		return res
	}

	res = res && (e.Size == uint64(len(e.Data32)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 32 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.Data32)))
		return res
	}

	for i := range (*d).Data32 {
		res = res && ((*d).Data32[i] == (*e).Data32[i])

		if !res {
			fmt.Printf("IsEqual :: Failed matrix %s - #%d check 32 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).Data32[i], (*e).Data32[i])
		}
	}

	return res
}

func (d *DistanceMatrix1Dg) check64(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.Data64)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 64 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.Data64)))
		return res
	}

	res = res && (e.Size == uint64(len(e.Data64)))

	if !res {
		fmt.Printf("IsEqual :: Failed matrix %s - #%d check 64 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.Data64)))
		return res
	}

	for i := range (*d).Data64 {
		res = res && ((*d).Data64[i] == (*e).Data64[i])

		if !res {
			fmt.Printf("IsEqual :: Failed matrix %s - #%d check 64 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).Data64[i], (*e).Data64[i])
		}
	}

	return res
}

//
// Get
//

func (d *DistanceMatrix1Dg) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	p := d.ijToK(p1, p2)

	if d.NumBits == 16 {
		return uint64((*d).Data16[p])
	} else if d.NumBits == 32 {
		return uint64((*d).Data32[p])
	} else if d.NumBits == 64 {
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

//
// Check
//
func (d *DistanceMatrix1Dg) Check() (res bool) {
	res = true

	return res
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
		fmt.Printf("saving matrix            :  %-70s block num: %d block pos: %d\n", outPrefix, d.BlockNumber, d.BlockPosition)
		saver.Save(d)
	} else {
		fmt.Printf("loading matrix           :  %-70s block num: %d block pos: %d\n", outPrefix, d.BlockNumber, d.BlockPosition)
		saver.Load(d)
	}
}
