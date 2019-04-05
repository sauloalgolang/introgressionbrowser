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
	BlockNumber uint64
	MinPosition uint64
	MaxPosition uint64
	NumSNPS     uint64
	NumSamples  uint64
	Matrix      [][]uint64
}

func NewIBBlock(blockNumber uint64, numSamples uint64) *IBBlock {
	ibb := IBBlock{
		BlockNumber: blockNumber,
		MinPosition: 0,
		MaxPosition: 0,
		NumSNPS:     0,
		NumSamples:  numSamples,
		Matrix:      make([][]uint64, numSamples*numSamples, numSamples*numSamples),
	}

	return &ibb
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

	ibb.NumSNPS++

	ibb.MinPosition = tools.Min64(ibb.MinPosition, reg.Pos)
	ibb.MaxPosition = tools.Max64(ibb.MaxPosition, reg.Pos)

	if false {
		fmt.Println("Failure getting block")
		os.Exit(1)
	}
}
