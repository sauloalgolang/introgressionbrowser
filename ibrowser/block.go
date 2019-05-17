package ibrowser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
)

//
//
// BLOCK SECTION
//
//

// IBBlock holds the information of a block
type IBBlock struct {
	ChromosomeName   string
	ChromosomeNumber int
	BlockSize        uint64
	CounterBits      uint64
	MinPosition      uint64
	MaxPosition      uint64
	NumSNPS          uint64
	NumSamples       uint64
	BlockPosition    uint64
	BlockNumber      uint64
	Serial           uint64
	Matrix           *IBDistanceMatrix
	hasSerial        bool
}

// NewIBBlock generates a new instance of a block
func NewIBBlock(
	chromosomeName string,
	chromosomeNumber int,
	blockSize uint64,
	counterBits uint64,
	numSamples uint64,
	blockNumber uint64,
	blockPosition uint64,
) *IBBlock {
	log.Println("   NewIBBlock :: chromosomeName: ", chromosomeName,
		" chromosomeNumber: ", chromosomeNumber,
		" blockSize: ", blockSize,
		" blockNumber: ", blockNumber,
		" blockPosition: ", blockPosition,
		" numSamples: ", numSamples,
	)

	ibb := IBBlock{
		ChromosomeName:   chromosomeName,
		ChromosomeNumber: chromosomeNumber,
		BlockSize:        blockSize,
		CounterBits:      counterBits,
		MinPosition:      math.MaxUint64,
		MaxPosition:      0,
		NumSNPS:          0,
		NumSamples:       numSamples,
		BlockNumber:      blockNumber,
		BlockPosition:    blockPosition,
		Serial:           0,
		hasSerial:        false,
		Matrix: NewDistanceMatrix(
			chromosomeName,
			blockSize,
			counterBits,
			numSamples,
			blockPosition,
			blockNumber,
		),
	}

	return &ibb
}

func (ibb *IBBlock) String() string {
	return fmt.Sprint("Block :: ",
		" ChromosomeName:   ", ibb.ChromosomeName, "\n",
		" ChromosomeNumber: ", ibb.ChromosomeNumber, "\n",
		" BlockSize:        ", ibb.BlockSize, "\n",
		" CounterBits:      ", ibb.CounterBits, "\n",
		" MinPosition:      ", ibb.MinPosition, "\n",
		" MaxPosition:      ", ibb.MaxPosition, "\n",
		" NumSNPS:          ", ibb.NumSNPS, "\n",
		" NumSamples:       ", ibb.NumSamples, "\n",
		" BlockPosition:    ", ibb.BlockPosition, "\n",
		" BlockNumber:      ", ibb.BlockNumber, "\n",
		" Serial:           ", ibb.Serial, "\n",
	)
}

// AddVcfMatrix add a vcf matrix to the blocks
func (ibb *IBBlock) AddVcfMatrix(position uint64, distance *VCFDistanceMatrix) {
	// log.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = Min64(ibb.MinPosition, position)
	ibb.MaxPosition = Max64(ibb.MaxPosition, position)
	ibb.Matrix.IncrementWithVcfMatrix(distance)
}

// Add add another ibmatrix to the blocks
func (ibb *IBBlock) Add(position uint64, distance *IBDistanceMatrix) {
	// log.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = Min64(ibb.MinPosition, position)
	ibb.MaxPosition = Max64(ibb.MaxPosition, position)
	ibb.Matrix.Merge(distance)
}

// GetMatrix gets the summary matrix
func (ibb *IBBlock) GetMatrix() (*IBDistanceMatrix, bool) {
	return ibb.Matrix, true
}

// GetMatrixTable gets the table from the summary matrix
func (ibb *IBBlock) GetMatrixTable() (*IBDistanceTable, bool) {
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		return nil, false
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, false
	}

	return table, true
}

// GetColumn returns a give column from the table
func (ibb *IBBlock) GetColumn(referenceNumber int) (*IBDistanceTable, bool) {
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		return nil, false
	}

	col, hasCol := matrix.GetColumn(referenceNumber)

	if !hasCol {
		return nil, false
	}

	return col, hasCol
}

// Sum returns the sum between another block and this block
func (ibb *IBBlock) Sum(other *IBBlock) {
	ibb.NumSNPS += other.NumSNPS
	ibb.MinPosition = Min64(ibb.MinPosition, other.MinPosition)
	ibb.MaxPosition = Max64(ibb.MaxPosition, other.MaxPosition)

	matrix, hasMatrix := other.GetMatrix()

	if !hasMatrix {
		panic("no matrix")
	}

	ibb.Matrix.Merge(matrix)
}

// IsEqual returns wether this block is equal to the other block
func (ibb *IBBlock) IsEqual(other *IBBlock) (res bool) {
	res = true

	res = res && (ibb.NumSNPS == other.NumSNPS)

	if !res {
		log.Printf("IsEqual :: Failed block %s - #%d check - NumSNPS: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.NumSNPS, other.NumSNPS)
		return res
	}

	if ibb.NumSNPS > 0 {
		res = res && (ibb.MinPosition == other.MinPosition)

		if !res {
			log.Printf("IsEqual :: Failed block %s - #%d check - MinPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, other.MinPosition)
			return res
		}

		res = res && (ibb.MaxPosition == other.MaxPosition)

		if !res {
			log.Printf("IsEqual :: Failed block %s - #%d check - MaxPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MaxPosition, other.MaxPosition)
			return res
		}
	}

	matrix, _ := other.GetMatrix()

	res = res && ibb.Matrix.IsEqual(matrix)

	if !res {
		log.Printf("IsEqual :: Failed block %s - #%d check - Matrix not equal\n", ibb.ChromosomeName, ibb.BlockNumber)
		return res
	}

	return res
}

//
// Serial
//

// SetSerial sets the serial number of this block
func (ibb *IBBlock) SetSerial(serial uint64) {
	ibb.Serial = serial
	ibb.Matrix.Serial = serial
}

// CheckSerial checks wether the serial numbers match
func (ibb *IBBlock) CheckSerial(serial uint64) bool {
	eq1 := ibb.Serial == serial

	if !eq1 {
		log.Println("block serial ", ibb.Serial, " != ", serial, ibb)
	} else {
		// log.Println("block serial ", ibb.Serial, " == ", serial, ibb)
	}

	eq2 := ibb.Matrix.Serial == serial

	if !eq2 {
		log.Println("matrix serial ", ibb.Matrix.Serial, " != ", serial, ibb.Matrix)
	} else {
		// log.Println("matrix serial ", ibb.Matrix.Serial, " == ", serial, ibb.Matrix)
	}

	return eq1 && eq2
}

//
// Check
//

// Check checks the self consistency of the data
func (ibb *IBBlock) Check() (res bool) {
	res = true

	res = res && (ibb.BlockNumber == ibb.Matrix.BlockNumber)

	if !res {
		log.Printf("Failed block %s - #%d check - BlockNumber: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockNumber, ibb.Matrix.BlockNumber)
		return res
	}

	res = res && (ibb.BlockPosition == ibb.Matrix.BlockPosition)

	if !res {
		log.Printf("Failed block %s - #%d check - BlockPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockPosition, ibb.Matrix.BlockPosition)
		return res
	}

	// res = res && (ibb.ChromosomeName == ibb.Matrix.ChromosomeName)

	// if !res {
	// 	log.Printf("Failed block %s - #%d check - ChromosomeName: %s != %s\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.ChromosomeName, ibb.Matrix.ChromosomeName)
	// 	return res
	// }

	if ibb.NumSNPS > 0 {
		res = res && (ibb.MinPosition <= ibb.MaxPosition)

		if !res {
			log.Printf("Failed block %s - #%d check - MinPosition %d > MaxPosition %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, ibb.MaxPosition)
			log.Println(ibb)
			return res
		}
	}

	return res
}

//
// Filename
//

// GenFilename generates the filename to save this block
func (ibb *IBBlock) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + fmt.Sprintf("%012d", ibb.BlockNumber)

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Dump
//

// Dump dumps the matrix table to a binary file
func (ibb *IBBlock) Dump(dumper *MultiArrayFile) (serial uint64) {
	serial = uint64(0)
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		log.Println("failed getting matrix")
		os.Exit(1)
	}

	serial = matrix.Dump(dumper)
	ibb.SetSerial(serial)

	return
}

// UnDump reads the matrix table from a binary file
func (ibb *IBBlock) UnDump(dumper *MultiArrayFile) (serial uint64, hasData bool) {
	// log.Println("UnDump matrix :: ", ibb)

	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		log.Println("block.UnDump failed getting matrix")
		os.Exit(1)
	}

	log.Println("block.UnDump reading from file")
	serial, hasData = matrix.UnDump(dumper)
	log.Println("block.UnDump reading from file - DONE")

	if !hasData {
		log.Println("block.UnDump Tried to read beyond the file")
		os.Exit(1)
	}

	if !ibb.CheckSerial(serial) {
		log.Println("block.UnDump Mismatch in order of files")
		os.Exit(1)
	}

	log.Println("UnDump matrix :: ", ibb.BlockNumber, " - DONE")

	return
}
