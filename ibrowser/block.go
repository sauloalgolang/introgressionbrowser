package ibrowser

import (
	"fmt"
	"math"
)

//
//
// BLOCK SECTION
//
//

type IBBlock struct {
	ChromosomeName string
	BlockSize      uint64
	CounterBits    int
	BlockPosition  uint64
	BlockNumber    uint64
	MinPosition    uint64
	MaxPosition    uint64
	NumSNPS        uint64
	NumSamples     uint64
	Serial         int64
	matrix         *DistanceMatrix
}

func NewIBBlock(chromosomeName string, blockSize uint64, counterBits int, numSamples uint64, blockPosition uint64, blockNumber uint64) *IBBlock {
	fmt.Println("   NewIBBlock :: chromosomeName: ", chromosomeName,
		" blockSize: ", blockSize,
		" blockPosition: ", blockPosition,
		" blockNumber: ", blockNumber,
		" numSamples: ", numSamples,
	)

	ibb := IBBlock{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		CounterBits:    counterBits,
		NumSamples:     numSamples,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		MinPosition:    math.MaxUint64,
		MaxPosition:    0,
		NumSNPS:        0,
		Serial:         -1,
		matrix: NewDistanceMatrix(
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

func (ibb *IBBlock) Add(position uint64, distance *DistanceMatrix) {
	// fmt.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = Min64(ibb.MinPosition, position)
	ibb.MaxPosition = Max64(ibb.MaxPosition, position)
	ibb.matrix.Add(distance)
}

func (ibb *IBBlock) GetMatrix() *DistanceMatrix {
	return ibb.matrix
}

func (ibb *IBBlock) Sum(other *IBBlock) {
	ibb.NumSNPS += other.NumSNPS
	ibb.MinPosition = Min64(ibb.MinPosition, other.MinPosition)
	ibb.MaxPosition = Max64(ibb.MaxPosition, other.MaxPosition)
	ibb.matrix.Add(other.GetMatrix())
}

func (ibb *IBBlock) IsEqual(other *IBBlock) (res bool) {
	res = true

	res = res && (ibb.NumSNPS == other.NumSNPS)

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - NumSNPS: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.NumSNPS, other.NumSNPS)
		return res
	}

	res = res && (ibb.MinPosition == other.MinPosition)

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - MinPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, other.MinPosition)
		return res
	}

	res = res && (ibb.MaxPosition == other.MaxPosition)

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - MaxPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MaxPosition, other.MaxPosition)
		return res
	}

	res = res && ibb.matrix.IsEqual(other.GetMatrix())

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - Matrix not equal\n", ibb.ChromosomeName, ibb.BlockNumber)
		return res
	}

	return res
}

//
// Check
//

func (ibb *IBBlock) Check() (res bool) {
	res = true

	res = res && (ibb.BlockNumber == ibb.matrix.BlockNumber)

	if !res {
		fmt.Printf("Failed block %s - #%d check - BlockNumber: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockNumber, ibb.matrix.BlockNumber)
		return res
	}

	res = res && (ibb.BlockPosition == ibb.matrix.BlockPosition)

	if !res {
		fmt.Printf("Failed block %s - #%d check - BlockPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockPosition, ibb.matrix.BlockPosition)
		return res
	}

	// res = res && (ibb.ChromosomeName == ibb.matrix.ChromosomeName)

	// if !res {
	// 	fmt.Printf("Failed block %s - #%d check - ChromosomeName: %s != %s\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.ChromosomeName, ibb.matrix.ChromosomeName)
	// 	return res
	// }

	res = res && (ibb.MinPosition <= ibb.MaxPosition)

	if !res {
		fmt.Printf("Failed block %s - #%d check - MinPosition %d > MaxPosition %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, ibb.MaxPosition)
		return res
	}

	return res
}

//
// Filename
//

func (ibb *IBBlock) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + fmt.Sprintf("%012d", ibb.BlockNumber)

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//

func (ibb *IBBlock) Save(outPrefix string, format string, compression string) {
	ibb.saveLoad(true, outPrefix, format, compression)
}

//
// Load
//

func (ibb *IBBlock) Load(outPrefix string, format string, compression string) {
	ibb.saveLoad(false, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ibb *IBBlock) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := ibb.GenFilename(outPrefix, format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Printf("saving block             :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Save(ibb)
		ibb.matrix.Save(baseName, format, compression)
	} else {
		fmt.Printf("loading block            :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Load(ibb)

		ibb.matrix = NewDistanceMatrix(
			ibb.ChromosomeName,
			ibb.BlockSize,
			ibb.CounterBits,
			ibb.NumSamples,
			ibb.BlockPosition,
			ibb.BlockNumber,
		)

		ibb.matrix.Load(baseName, format, compression)
	}
}
