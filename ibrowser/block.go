package ibrowser

import (
	// "fmt"
	"math"
	// "os"
	"sync/atomic"
)

import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// BLOCK SECTION
//
//

type IBBlock struct {
	BlockNumber uint64
	MinPosition uint64
	MaxPosition uint64
	NumSNPS     uint64
	NumSamples  uint64
	Matrix      *tools.DistanceMatrix
}

func NewIBBlock(blockNumber uint64, numSamples uint64) *IBBlock {
	ibb := IBBlock{
		BlockNumber: blockNumber,
		MinPosition: math.MaxUint64,
		MaxPosition: 0,
		NumSNPS:     0,
		NumSamples:  numSamples,
		Matrix:      tools.NewDistanceMatrix(numSamples),
	}

	return &ibb
}

func (ibb *IBBlock) add(position uint64, distance *tools.DistanceMatrix, isAtomic bool) {
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
		ibb.Matrix.AddAtomic(distance)
	} else {
		ibb.Matrix.Add(distance)
	}
}

func (ibb *IBBlock) Add(position uint64, distance *tools.DistanceMatrix) {
	ibb.add(position, distance, false)
}

func (ibb *IBBlock) AddAtomic(position uint64, distance *tools.DistanceMatrix) {
	ibb.add(position, distance, true)
}
