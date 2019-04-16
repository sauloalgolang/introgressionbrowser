package ibrowser

import (
	"fmt"
	"math"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

//
//
// BLOCK SECTION
//
//

type IBBlock struct {
	ChromosomeName string
	BlockSize      uint64
	NumBits        int
	BlockPosition  uint64
	BlockNumber    uint64
	MinPosition    uint64
	MaxPosition    uint64
	NumSNPS        uint64
	NumSamples     uint64
	matrix         *interfaces.DistanceMatrix
}

func NewIBBlock(chromosomeName string, blockSize uint64, numBits int, numSamples uint64, blockPosition uint64, blockNumber uint64) *IBBlock {
	fmt.Println("   NewIBBlock :: chromosomeName: ", chromosomeName,
		" blockSize: ", blockSize,
		" blockPosition: ", blockPosition,
		" blockNumber: ", blockNumber,
		" numSamples: ", numSamples,
	)

	ibb := IBBlock{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		NumBits:        numBits,
		NumSamples:     numSamples,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		MinPosition:    math.MaxUint64,
		MaxPosition:    0,
		NumSNPS:        0,
		matrix: interfaces.NewDistanceMatrix(
			chromosomeName,
			blockSize,
			numBits,
			numSamples,
			blockPosition,
			blockNumber,
		),
	}

	return &ibb
}

func (ibb *IBBlock) Add(position uint64, distance *interfaces.DistanceMatrix) {
	// fmt.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = tools.Min64(ibb.MinPosition, position)
	ibb.MaxPosition = tools.Max64(ibb.MaxPosition, position)
	ibb.matrix.Add(distance)
}

func (ibb *IBBlock) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + fmt.Sprintf("%012d", ibb.BlockNumber)

	saver := save.NewSaverCompressed(baseName, format, compression)

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
	saver := save.NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Printf("saving block             :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Save(ibb)
		ibb.matrix.Save(baseName, format, compression)
	} else {
		fmt.Printf("loading block            :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Load(ibb)

		ibb.matrix = interfaces.NewDistanceMatrix(
			ibb.ChromosomeName,
			ibb.BlockSize,
			ibb.NumBits,
			ibb.NumSamples,
			ibb.BlockPosition,
			ibb.BlockNumber,
		)

		ibb.matrix.Load(baseName, format, compression)
	}
}
