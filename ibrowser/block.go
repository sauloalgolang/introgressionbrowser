package ibrowser

import (
	"fmt"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"
import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// BLOCK SECTION
//
//

type IBBlock struct {
	blockNumber uint64
	minPosition uint64
	maxPosition uint64
	numSNPS     uint64
	numSamples  uint64
	matrix      [][]uint64
}

func NewIBBlock(blockNumber uint64, numSamples uint64) IBBlock {
	ibb := IBBlock{
		blockNumber: blockNumber,
		minPosition: 0,
		maxPosition: 0,
		numSNPS:     0,
		numSamples:  numSamples,
		matrix:      make([][]uint64, numSamples*numSamples, numSamples*numSamples),
	}

	return ibb
}

func (ibb *IBBlock) Add(reg *interfaces.VCFRegister) {
	// type VCFRegister struct {
	//  Samples      *[]string
	// 	IsHomozygous bool
	// 	IsIndel      bool
	// 	IsMNP        bool
	//  Row          uint64
	// 	Chromosome   string
	// 	Position     uint64
	// 	Quality      float32
	// 	Info         map[string]interface{}
	// 	Filter       string
	// 	NumAlt       uint64
	// 	Phased       bool
	// 	GT           [][]int
	// 	Fields       map[string]string
	// }

	ibb.numSNPS++

	ibb.minPosition = tools.Min64(ibb.minPosition, reg.Position)
	ibb.maxPosition = tools.Max64(ibb.maxPosition, reg.Position)

	if false {
		fmt.Println("Failure getting block")
		os.Exit(1)
	}
}
