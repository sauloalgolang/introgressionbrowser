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
	// type Variant struct {
	// 	Chromosome      string
	// 	Pos        		uint64
	// 	Id         		string
	// 	Ref        		string
	// 	Alt        		[]string
	// 	Quality    		float32
	// 	Filter     		string
	// 	Info       		InfoMap
	// 	Format     		[]string
	// 	Samples    		[]*SampleGenotype
	// 	Header     		*Header
	// 	LineNumber 		int64
	// }

	ibb.numSNPS++

	ibb.minPosition = tools.Min64(ibb.minPosition, reg.Pos)
	ibb.maxPosition = tools.Max64(ibb.maxPosition, reg.Pos)

	if false {
		fmt.Println("Failure getting block")
		os.Exit(1)
	}
}
