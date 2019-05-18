package ibrowser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

//
//
// Matrix 1D
//
//

// DistanceMatrix1Dg is the definitions of a 1d matrix
type DistanceMatrix1Dg struct {
	ChromosomeName string
	BlockSize      uint64
	Dimension      uint64
	Size           uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Serial         uint64
	CounterBits    uint64
	data16         DistanceRow16
	data32         DistanceRow32
	data64         DistanceRow64
	tm             *TriangularMatrix
	// Data           []interface{}
}

func (d *DistanceMatrix1Dg) String() string {
	return fmt.Sprint("Matrix :: ",
		" ChromosomeName: ", d.ChromosomeName, "\n",
		" BlockSize:      ", d.BlockSize, "\n",
		" Dimension:      ", d.Dimension, "\n",
		" Size:           ", d.Size, "\n",
		" BlockPosition:  ", d.BlockPosition, "\n",
		" BlockNumber:    ", d.BlockNumber, "\n",
		" Serial:         ", d.Serial, "\n",
		" CounterBits:    ", d.CounterBits, "\n",
	)
}

// NewDistanceMatrix1Dg16 creates a new instance of 1d matrix using 16 bit numbers
func NewDistanceMatrix1Dg16(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 16, dimension, blockPosition, blockNumber)
}

// NewDistanceMatrix1Dg32 creates a new instance of 1d matrix using 32 bit numbers
func NewDistanceMatrix1Dg32(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 32, dimension, blockPosition, blockNumber)
}

// NewDistanceMatrix1Dg64 creates a new instance of 1d matrix using 64 bit numbers
func NewDistanceMatrix1Dg64(chromosomeName string, blockSize uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	return NewDistanceMatrix1Dg(chromosomeName, blockSize, 64, dimension, blockPosition, blockNumber)
}

// NewDistanceMatrix1Dg creates a new instance of 1d matrix using counterBits bit numbers
func NewDistanceMatrix1Dg(chromosomeName string, blockSize uint64, counterBits uint64, dimension uint64, blockPosition uint64, blockNumber uint64) *DistanceMatrix1Dg {
	log.Println("    NewDistanceMatrix1D :: Chromosome: ", chromosomeName,
		" Block Size: ", blockSize,
		" Bits:", counterBits,
		" Dimension:", dimension,
		" Block Position: ", blockPosition,
		" Block Number: ", blockNumber,
	)

	d := DistanceMatrix1Dg{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		Dimension:      dimension,
		CounterBits:    counterBits,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		Serial:         0,
	}

	d.init()

	if d.CounterBits == 16 {
		// d.Data = make(DistanceRow32, d.Size, d.Size)
		d.data16 = make(DistanceRow16, d.Size, d.Size)
		d.data32 = make(DistanceRow32, 0, 0)
		d.data64 = make(DistanceRow64, 0, 0)
	} else if d.CounterBits == 32 {
		// d.Data = make(DistanceRow32, d.Size, d.Size)
		d.data16 = make(DistanceRow16, 0, 0)
		d.data32 = make(DistanceRow32, d.Size, d.Size)
		d.data64 = make(DistanceRow64, 0, 0)
	} else if d.CounterBits == 64 {
		// d.Data = make(DistanceRow64, d.Size, d.Size)
		d.data16 = make(DistanceRow16, 0, 0)
		d.data32 = make(DistanceRow32, 0, 0)
		d.data64 = make(DistanceRow64, d.Size, d.Size)
	}

	d.Clean()

	return &d
}

func (d *DistanceMatrix1Dg) init() {
	d.tm = NewTriangularMatrix(d.Dimension)
	d.updateSize()
	// testIJP()
}

//
// GetTable
//

// GetTable returns data table
func (d *DistanceMatrix1Dg) GetTable() (*DistanceRow64, bool) {
	if d.CounterBits == 64 {
		return &d.data64, true
	}

	data := make(DistanceRow64, d.Size, d.Size)

	if d.CounterBits == 16 {
		for i := range (*d).data16 {
			data[i] = uint64((*d).data16[i])
		}
		return &d.data64, true
	} else if d.CounterBits == 32 {
		for i := range (*d).data32 {
			data[i] = uint64((*d).data32[i])
		}
		return &d.data64, true
	}

	return nil, false
}

// GetTable16 return data table using 16 bits
func (d *DistanceMatrix1Dg) GetTable16() *DistanceRow16 {
	if d.CounterBits != 16 {
		log.Println("calling GetTable16 when numbits not 16")
		os.Exit(1)
	}

	return &d.data16
}

// GetTable32 return data table using 32 bits
func (d *DistanceMatrix1Dg) GetTable32() *DistanceRow32 {
	if d.CounterBits != 32 {
		log.Println("calling GetTable32 when numbits not 32")
		os.Exit(1)
	}

	return &d.data32
}

// GetTable64 return data table using 64 bits
func (d *DistanceMatrix1Dg) GetTable64() *DistanceRow64 {
	if d.CounterBits != 64 {
		log.Println("calling GetTable64 when numbits not 64")
		os.Exit(1)
	}

	return &d.data64
}

//
// Get Column

// GetColumn gets a single column from the matrix
func (d *DistanceMatrix1Dg) GetColumn(columNumber int) (*DistanceRow64, bool) {
	dr := make(DistanceRow64, d.Dimension, d.Dimension)

	for p := uint64(0); p < uint64(columNumber); p++ {
		dr[p] = d.GetPos(uint64(columNumber), p)
	}

	return &dr, true
}

//
// Clean
//

// Clean zeroes the matrix
func (d *DistanceMatrix1Dg) Clean() {
	if d.CounterBits == 16 {
		d.clean16()
	} else if d.CounterBits == 32 {
		d.clean32()
	} else if d.CounterBits == 64 {
		d.clean64()
	}
}

// Clean16 zeroes the 16 bit matrix
func (d *DistanceMatrix1Dg) clean16() {
	for i := range (*d).data16 {
		(*d).data16[i] = uint16(0)
	}
}

// Clean32 zeroes the 32 bit matrix
func (d *DistanceMatrix1Dg) clean32() {
	for i := range (*d).data32 {
		(*d).data32[i] = uint32(0)
	}
}

// Clean64 zeroes the 64 bit matrix
func (d *DistanceMatrix1Dg) clean64() {
	for i := range (*d).data64 {
		(*d).data64[i] = uint64(0)
	}
}

//
// Increment
//

// IncrementWithVcfMatrix increments the table using a vcf matrix
func (d *DistanceMatrix1Dg) IncrementWithVcfMatrix(e *VCFDistanceMatrix) {
	j := uint64(0)
	v := uint64(0)
	le := uint64(len(*e))
	for i := uint64(1); i < le; i++ {
		for j = i; j < le; j++ {
			v = (*e)[i][j]
			d.Increment(i, j, v)
		}
	}
}

// Increment increments the table using another table
func (d *DistanceMatrix1Dg) Increment(p1 uint64, p2 uint64, val uint64) {
	p := d.ijToP(p1, p2)

	if d.CounterBits == 16 {
		d.increment16(p, val)
	} else if d.CounterBits == 32 {
		d.increment32(p, val)
	} else if d.CounterBits == 64 {
		d.increment64(p, val)
	}
}

// Increment16 increments the 16 bits table using a vcf matrix
func (d *DistanceMatrix1Dg) increment16(p uint64, val uint64) {
	if val >= uint64(math.MaxUint16) {
		log.Println("count 16 overflow")
		os.Exit(1)
	}

	v := uint64((*d).data16[p])
	r := v + val

	(*d).data16[p] = uint16(r)
}

// Increment32 increments the 32 bits table using a vcf matrix
func (d *DistanceMatrix1Dg) increment32(p uint64, val uint64) {
	if val >= uint64(math.MaxUint32) {
		log.Println("count 32 overflow")
		os.Exit(1)
	}

	v := uint64((*d).data32[p])
	r := v + val

	if r >= uint64(math.MaxUint32) {
		log.Println("count 32 overflow")
		os.Exit(1)
	}

	(*d).data32[p] = uint32(r)
}

// Increment64 increments the 64 bits table using a vcf matrix
func (d *DistanceMatrix1Dg) increment64(p uint64, val uint64) {
	(*d).data64[p] += val
}

//
// Set
//

// Set sets the value for matrix position p1,p2
func (d *DistanceMatrix1Dg) Set(p1 uint64, p2 uint64, val uint64) {
	p := d.ijToP(p1, p2)

	if d.CounterBits == 16 {
		d.set16(p, val)
	} else if d.CounterBits == 32 {
		d.set32(p, val)
	} else if d.CounterBits == 64 {
		d.set64(p, val)
	}
}

// Set16 sets the value for matrix position p1,p2 of 16 bit matrix
func (d *DistanceMatrix1Dg) set16(p uint64, val uint64) {
	if val >= uint64(math.MaxUint16) {
		log.Println("count 16 overflow")
		os.Exit(1)
	}

	(*d).data16[p] = uint16(val)
}

// Set32 sets the value for matrix position p1,p2 of 32 bit matrix
func (d *DistanceMatrix1Dg) set32(p uint64, val uint64) {
	if val >= uint64(math.MaxUint32) {
		log.Println("count 32 overflow")
		os.Exit(1)
	}

	(*d).data32[p] = uint32(val)
}

// Set64 sets the value for matrix position p1,p2 of 64 bit matrix
func (d *DistanceMatrix1Dg) set64(p uint64, val uint64) {
	(*d).data64[p] = val
}

//
// Merge
//

// Merge merges two distance matrices
func (d *DistanceMatrix1Dg) Merge(e *DistanceMatrix1Dg) {
	if d.CounterBits == 16 {
		d.merge16(e)
	} else if d.CounterBits == 32 {
		d.merge32(e)
	} else if d.CounterBits == 64 {
		d.merge64(e)
	}
}

// merge16 merges a distance matrices on the 16 bit table
func (d *DistanceMatrix1Dg) merge16(e *DistanceMatrix1Dg) {
	mi := uint64(math.MaxInt16)
	for i := range (*d).data16 {
		if uint64((*d).data16[i])+uint64((*e).data16[i]) >= mi {
			log.Println("counter 16 overflow")
			os.Exit(1)
		}
		(*d).data16[i] += (*e).data16[i]
	}
}

// merge32 merges a distance matrices on the 32 bit table
func (d *DistanceMatrix1Dg) merge32(e *DistanceMatrix1Dg) {
	mi := uint64(math.MaxInt32)
	for i := range (*d).data32 {
		vdi := uint64((*d).data32[i])
		vei := uint64((*e).data32[i])
		if (vdi + vei) >= mi {
			log.Println("counter 32 overflow", vdi, vei, mi)
			os.Exit(1)
		}
		(*d).data32[i] += (*e).data32[i]
	}
}

// merge64 merges a distance matrices on the 64 bit table
func (d *DistanceMatrix1Dg) merge64(e *DistanceMatrix1Dg) {
	for i := range (*d).data64 {
		(*d).data64[i] += (*e).data64[i]
	}
}

//
// IsEqual
//

// IsEqual checks whether two tables are equal
func (d *DistanceMatrix1Dg) IsEqual(e *DistanceMatrix1Dg) (res bool) {
	res = true

	// res = res && (d.ChromosomeName == e.ChromosomeName)
	//
	// if !res {
	// 	log.Printf("IsEqual :: Failed matrix %s - #%d check - ChromosomeName %s != %s\n", d.ChromosomeName, d.BlockNumber, d.ChromosomeName, e.ChromosomeName)
	// 	return res
	// }

	res = res && (d.BlockSize == e.BlockSize)

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check - BlockSize %d != %d\n", d.ChromosomeName, d.BlockNumber, d.BlockSize, e.BlockSize)
		return res
	}

	res = res && (d.Dimension == e.Dimension)

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check - Dimension %d != %d\n", d.ChromosomeName, d.BlockNumber, d.Dimension, e.Dimension)
		return res
	}

	res = res && (d.CounterBits == e.CounterBits)

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check - CounterBits %d != %d\n", d.ChromosomeName, d.BlockNumber, d.CounterBits, e.CounterBits)
		return res
	}

	res = res && (d.Size == e.Size)

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check - Size %d != %d\n", d.ChromosomeName, d.BlockNumber, d.Size, e.Size)
		return res
	}

	if d.CounterBits == 16 {
		d.isEqual16(e)
	} else if d.CounterBits == 32 {
		d.isEqual32(e)
	} else if d.CounterBits == 64 {
		d.isEqual64(e)
	}

	return res

}

func (d *DistanceMatrix1Dg) isEqual16(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.data16)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 16 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.data16)))
		return res
	}

	res = res && (e.Size == uint64(len(e.data16)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 16 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.data16)))
		return res
	}

	for i := range (*d).data16 {
		res = res && ((*d).data16[i] == (*e).data16[i])

		if !res {
			log.Printf("IsEqual :: Failed matrix %s - #%d check 16 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).data16[i], (*e).data16[i])
		}
	}

	return res
}

func (d *DistanceMatrix1Dg) isEqual32(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.data32)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 32 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.data32)))
		return res
	}

	res = res && (e.Size == uint64(len(e.data32)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 32 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.data32)))
		return res
	}

	for i := range (*d).data32 {
		res = res && ((*d).data32[i] == (*e).data32[i])

		if !res {
			log.Printf("IsEqual :: Failed matrix %s - #%d check 32 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).data32[i], (*e).data32[i])
		}
	}

	return res
}

func (d *DistanceMatrix1Dg) isEqual64(e *DistanceMatrix1Dg) (res bool) {
	res = true

	res = res && (d.Size == uint64(len(d.data64)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 64 - D Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, d.Size, uint64(len(d.data64)))
		return res
	}

	res = res && (e.Size == uint64(len(e.data64)))

	if !res {
		log.Printf("IsEqual :: Failed matrix %s - #%d check 64 - E Size %d != Len %d\n", d.ChromosomeName, d.BlockNumber, e.Size, uint64(len(e.data64)))
		return res
	}

	for i := range (*d).data64 {
		res = res && ((*d).data64[i] == (*e).data64[i])

		if !res {
			log.Printf("IsEqual :: Failed matrix %s - #%d check 64 - Position %d : %d != %d\n", d.ChromosomeName, d.BlockNumber, i, (*d).data64[i], (*e).data64[i])
		}
	}

	return res
}

//
// Get
//

// GetPos returns the value at a given position in the table
func (d *DistanceMatrix1Dg) GetPos(p1 uint64, p2 uint64) uint64 {
	p := d.ijToP(p1, p2)

	log.Printf("GetPos :: p1 %d p2 %d p %d len data16 %d data32 %d data64 %d", p1, p2, p, len((*d).data16), len((*d).data32), len((*d).data64))

	if d.CounterBits == 16 {
		return uint64((*d).data16[p])
	} else if d.CounterBits == 32 {
		return uint64((*d).data32[p])
	} else if d.CounterBits == 64 {
		return (*d).data64[p]
	}

	return 0
}

// GenFilename generates the filename to save this matrix to file
func (d *DistanceMatrix1Dg) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "_matrix"

	saver := save.NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Unexported Methods
//

func (d *DistanceMatrix1Dg) ijToP(i uint64, j uint64) uint64 {
	return d.tm.IJToP(i, j)
}

func (d *DistanceMatrix1Dg) pToIJ(k uint64) (uint64, uint64) {
	return d.tm.PToIJ(k)
}

func (d *DistanceMatrix1Dg) calculateSize() (size uint64) {
	size = d.tm.CalculateSize()
	return
}

func (d *DistanceMatrix1Dg) updateSize() {
	d.Size = d.calculateSize()
}

//
// Check
//

// Check checks the self consistency of the data
func (d *DistanceMatrix1Dg) Check() (res bool) {
	res = true

	// *(data.(*[]uint32)) = e

	// d := DistanceMatrix1Dg{
	// 	ChromosomeName: chromosomeName,
	// 	BlockSize:      blockSize,
	// 	Dimension:      dimension,
	// 	Size:           size,
	// 	CounterBits:    counterBits,
	// 	BlockPosition:  blockPosition,
	// 	BlockNumber:    blockNumber,
	// 	Serial:         0,
	// }

	return res
}

func (d *DistanceMatrix1Dg) checkHeader(header RegisterHeader) {
	log.Println("DistanceMatrix1Dg::checkHeader")

	if header.Serial != d.Serial {
		log.Printf("DistanceMatrix1Dg:checkHeader :: Serial mismatch: %d != %d", header.Serial, d.Serial)
		panic("DistanceMatrix1Dg:checkHeader :: Serial mismatch")
	} else if header.CounterBits != d.CounterBits {
		log.Printf("DistanceMatrix1Dg:checkHeader :: CounterBits mismatch: %d != %d", header.CounterBits, d.CounterBits)
		panic("DistanceMatrix1Dg:checkHeader :: CounterBits mismatch")
	} else if len(d.data32) == 0 {
		log.Printf("DistanceMatrix1Dg:checkHeader :: Empty table: %d != %d", header.CounterBits, d.data32)
		panic("DistanceMatrix1Dg:checkHeader :: Empty table")
	} else if uint64(len(d.data32)) != d.Size {
		log.Printf("DistanceMatrix1Dg:checkHeader :: LenData mismatch: %d != %d", len(d.data32), d.Size)
		panic("DistanceMatrix1Dg:checkHeader :: LenData mismatch")
	} else if header.DataLen != d.Size {
		log.Printf("DistanceMatrix1Dg:checkHeader :: DataLen mismatch: %d != %d", header.DataLen, d.Size)
		panic("DistanceMatrix1Dg:checkHeader :: DataLen mismatch")
	} else {
		sumDataV := uint64(0)
		
		for _, w := range d.data32 {
			sumDataV += uint64(w)
		}

		if header.SumData != sumDataV {
			log.Printf("DistanceMatrix1Dg:checkHeader :: SumData mismatch: %d != %d", header.SumData, sumDataV)
			panic("DistanceMatrix1Dg:checkHeader :: SumData mismatch")
		}
	}

	log.Println("DistanceMatrix1Dg::checkHeader - DONE")
}


//
// Save and Load
//

// Save saves this matrix to file
func (d *DistanceMatrix1Dg) Save(outPrefix string, format string, compression string) {
	d.saveLoad(true, outPrefix, format, compression)
}

// Load loads a matrix from file
func (d *DistanceMatrix1Dg) Load(outPrefix string, format string, compression string) {
	d.saveLoad(false, outPrefix, format, compression)
}

func (d *DistanceMatrix1Dg) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := d.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)

	if isSave {
		log.Printf("saving matrix            :  %-70s block num: %d block pos: %d\n", outPrefix, d.BlockNumber, d.BlockPosition)
		saver.Save(d)
	} else {
		log.Printf("loading matrix           :  %-70s block num: %d block pos: %d\n", outPrefix, d.BlockNumber, d.BlockPosition)
		saver.Load(d)
		d.init()
	}
}

//
// Dump
//

// Dump dumps the table to binary file
func (d *DistanceMatrix1Dg) Dump(dumper *MultiArrayFile) (serial uint64) {
	serial = uint64(0)

	if d.CounterBits == 16 {
		serial = dumper.Write16(&d.data16)
	} else if d.CounterBits == 32 {
		serial = dumper.Write32(&d.data32)
	} else if d.CounterBits == 64 {
		serial = dumper.Write64(&d.data64)
	}

	return
}

// UnDump loads the table from a binary file
func (d *DistanceMatrix1Dg) UnDump(dumper *MultiArrayFile) (serial uint64, hasData bool) {
	header := RegisterHeader{}

	if d.CounterBits == 16 {
		header = dumper.Read16(&d.data16)
	} else if d.CounterBits == 32 {
		header = dumper.Read32(&d.data32)
	} else if d.CounterBits == 64 {
		header = dumper.Read64(&d.data64)
	}

	d.checkHeader(header)

	return header.Serial, header.HasData
}


func testIJP() {
	log.Println("testIJP")
	dimension := uint64(10)
	tm:= NewTriangularMatrix(dimension)
	size := tm.CalculateSize()

	var i, j, p, q uint64

	for p = 0; p < size; p++ {
		i, j = tm.PToIJ(p)
		q = tm.IJToP(i, j)
		log.Printf("d %3d size %3d p %3d i %3d j %3d q %3d", dimension, size, p, i, j, q)
		if p != q {
			log.Panic("Error")
		}
	}
}
