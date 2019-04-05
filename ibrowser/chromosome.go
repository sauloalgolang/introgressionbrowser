package ibrowser

import (
	"fmt"
	"os"
)

import "runtime/debug"

import "github.com/sauloalgolang/introgressionbrowser/interfaces"
import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	Chromosome     string
	MinPosition    uint64
	MaxPosition    uint64
	NumBlocks      uint64
	NumSNPS        uint64
	NumSamples     uint64
	KeepEmptyBlock bool
	Blocks         []*IBBlock
	BlockNames     map[uint64]uint64
}

func NewIBChromosome(chromosome string, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	ibc := IBChromosome{
		Chromosome:     chromosome,
		NumSamples:     numSamples,
		MinPosition:    0,
		MaxPosition:    0,
		NumBlocks:      0,
		NumSNPS:        0,
		KeepEmptyBlock: keepEmptyBlock,
		Blocks:         make([]*IBBlock, 0, 100),
		BlockNames:     make(map[uint64]uint64, 100),
	}

	return &ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) {
	ibc.Blocks = append(ibc.Blocks, NewIBBlock(blockNum, ibc.NumSamples))
	ibc.BlockNames[blockNum] = uint64(len(ibc.Blocks)) - uint64(1)
	ibc.NumBlocks++
}

func (ibc *IBChromosome) GetBlock(blockNum uint64) (*IBBlock, bool) {
	if blockPos, ok := ibc.BlockNames[blockNum]; ok {
		if blockPos >= uint64(len(ibc.Blocks)) {
			fmt.Println(&ibc, "Index out of range. block num:", blockNum, "block pos:", blockPos, "len:", len(ibc.Blocks), "NumBlocks:", ibc.NumBlocks)
			fmt.Println(&ibc, "BlockNames", ibc.BlockNames)
			fmt.Println(&ibc, "Blocks", ibc.Blocks)
			debug.PrintStack()
			os.Exit(1)
		}

		return ibc.Blocks[blockPos], ok
	} else {
		return &IBBlock{}, ok
	}
}

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) {
	if _, hasBlock := ibc.GetBlock(blockNum); !hasBlock {
		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(len(ibc.Blocks))

			if NumBlocks != 0 {
				lastBlockPos = NumBlocks - 1
			}

			for currBlockPos := lastBlockPos; currBlockPos < blockNum; currBlockPos++ {
				ibc.AppendBlock(currBlockPos)
			}
		}
		ibc.AppendBlock(blockNum)
	}
}

func (ibc *IBChromosome) Add(blockNum uint64, reg *interfaces.VCFRegister) {
	ibc.normalizeBlocks(blockNum)

	if block, success := ibc.GetBlock(blockNum); success {
		block.Add(reg)
		ibc.NumSNPS++
		ibc.MinPosition = tools.Min64(ibc.MinPosition, block.MinPosition)
		ibc.MaxPosition = tools.Max64(ibc.MaxPosition, block.MaxPosition)
	} else {
		fmt.Println("Failure getting block", blockNum)
		os.Exit(1)
	}
}
