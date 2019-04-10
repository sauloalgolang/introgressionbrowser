package ibrowser

import (
	"fmt"
	"math"
	// "os"
	"sync/atomic"
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
	BlockPosition  uint64
	BlockNumber    uint64
	MinPosition    uint64
	MaxPosition    uint64
	NumSNPS        uint64
	NumSamples     uint64
	matrix         *interfaces.DistanceMatrix
}

func NewIBBlock(chromosomeName string, blockSize uint64, blockPosition uint64, blockNumber uint64, numSamples uint64) *IBBlock {
	fmt.Println("NewIBBlock :: chromosomeName: ", chromosomeName, " blockSize: ", blockSize, "blockPosition: ", blockPosition, "blockNumber: ", blockNumber, " numSamples: ", numSamples)

	ibb := IBBlock{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		BlockPosition:  blockPosition,
		BlockNumber:    blockNumber,
		MinPosition:    math.MaxUint64,
		MaxPosition:    0,
		NumSNPS:        0,
		NumSamples:     numSamples,
		matrix:         interfaces.NewDistanceMatrix(chromosomeName, blockSize, blockPosition, blockNumber, numSamples),
	}

	return &ibb
}

func (ibb *IBBlock) add(position uint64, distance *interfaces.DistanceMatrix, isAtomic bool) {
	if isAtomic {
		atomic.AddUint64(&ibb.NumSNPS, 1)
	} else {
		ibb.NumSNPS++
	}

	if isAtomic {
		atomic.StoreUint64(&ibb.MinPosition, tools.Min64(atomic.LoadUint64(&ibb.MinPosition), position))
		atomic.StoreUint64(&ibb.MaxPosition, tools.Max64(atomic.LoadUint64(&ibb.MaxPosition), position))
	} else {
		ibb.MinPosition = tools.Min64(ibb.MinPosition, position)
		ibb.MaxPosition = tools.Max64(ibb.MaxPosition, position)
	}

	if isAtomic {
		ibb.matrix.AddAtomic(distance)
	} else {
		ibb.matrix.Add(distance)
	}
}

func (ibb *IBBlock) Add(position uint64, distance *interfaces.DistanceMatrix) {
	ibb.add(position, distance, false)
}

func (ibb *IBBlock) AddAtomic(position uint64, distance *interfaces.DistanceMatrix) {
	ibb.add(position, distance, true)
}

func (ibb *IBBlock) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix + "." + fmt.Sprintf("%012d", ibb.BlockNumber)
	fileName = save.GenFilename(baseName, format)
	return baseName, fileName
}

func (ibb *IBBlock) Save(outPrefix string, format string) {
	baseName, _ := ibb.GenFilename(outPrefix, format)
	save.Save(baseName, format, ibb)
	ibb.matrix.Save(baseName, format)
}
