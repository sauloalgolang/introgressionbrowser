package ibrowser

import (
	"fmt"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"
import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	chromosome     string
	minPosition    uint64
	maxPosition    uint64
	numBlocks      uint64
	numSNPS        uint64
	numSamples     uint64
	keepEmptyBlock bool
	blocks         []IBBlock
	blockNames     map[uint64]uint64
}

func NewIBChromosome(chromosome string, numSamples uint64, keepEmptyBlock bool) IBChromosome {
	ibc := IBChromosome{
		chromosome:     chromosome,
		numSamples:     numSamples,
		minPosition:    0,
		maxPosition:    0,
		numBlocks:      0,
		numSNPS:        0,
		keepEmptyBlock: keepEmptyBlock,
		blocks:         make([]IBBlock, 0, 100),
		blockNames:     make(map[uint64]uint64, 100),
	}

	return ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) {
	ibc.blocks = append(ibc.blocks, NewIBBlock(blockNum, ibc.numSamples))
	ibc.blockNames[blockNum] = uint64(len(ibc.blocks)) - uint64(1)
	ibc.numBlocks++
}

func (ibc *IBChromosome) HasBlock(blockNum uint64) bool {
	if _, ok := ibc.blockNames[blockNum]; ok {
		return ok
	} else {
		return ok
	}
}

func (ibc *IBChromosome) GetBlock(blockNum uint64) (IBBlock, bool) {
	if blockPos, ok := ibc.blockNames[blockNum]; ok {
		return ibc.blocks[blockPos], ok
	} else {
		return IBBlock{}, ok
	}
}

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) {
	if !ibc.HasBlock(blockNum) {
		if ibc.keepEmptyBlock {
			lastBlockPos := uint64(0)
			numBlocks := uint64(len(ibc.blocks))

			if numBlocks != 0 {
				lastBlockPos = numBlocks - 1
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
		ibc.numSNPS++
		ibc.minPosition = tools.Min64(ibc.minPosition, block.minPosition)
		ibc.maxPosition = tools.Max64(ibc.maxPosition, block.maxPosition)
	} else {
		fmt.Println("Failure getting block", blockNum)
		os.Exit(1)
	}
}
