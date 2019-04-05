package ibrowser

import (
	"fmt"
	"os"
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
	Matrix      tools.DistanceMatrix
}

func NewIBBlock(blockNumber uint64, numSamples uint64) *IBBlock {
	ibb := IBBlock{
		BlockNumber: blockNumber,
		MinPosition: 0,
		MaxPosition: 0,
		NumSNPS:     0,
		NumSamples:  numSamples,
		Matrix:      *tools.NewDistanceMatrix(numSamples),
	}

	return &ibb
}

func (ibb *IBBlock) Add(position uint64, distance *tools.DistanceMatrix) {
	ibb.NumSNPS++

	ibb.MinPosition = tools.Min64(ibb.MinPosition, position)
	ibb.MaxPosition = tools.Max64(ibb.MaxPosition, position)

	ibb.Matrix.Add(distance)

	if false {
		fmt.Println("Failure getting block")
		os.Exit(1)
	}
}
