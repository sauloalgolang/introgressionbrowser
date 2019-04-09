package ibrowser

import (
	"fmt"
	"math"
	"os"
)

import "runtime/debug"

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	ChromosomeName string
	BlockSize      uint64
	MinPosition    uint64
	MaxPosition    uint64
	NumBlocks      uint64
	NumSNPS        uint64
	NumSamples     uint64
	KeepEmptyBlock bool
	Block          *IBBlock
	BlockNames     map[uint64]uint64
	blocks         []*IBBlock
}

func NewIBChromosome(chromosomeName string, blockSize uint64, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	ibc := IBChromosome{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		NumSamples:     numSamples,
		MinPosition:    math.MaxUint64,
		MaxPosition:    0,
		NumBlocks:      0,
		NumSNPS:        0,
		KeepEmptyBlock: keepEmptyBlock,
		Block:          NewIBBlock(chromosomeName, blockSize, 0, 0, numSamples),
		BlockNames:     make(map[uint64]uint64, 100),
		blocks:         make([]*IBBlock, 0, 100),
	}

	return &ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) {
	blockPos := uint64(len(ibc.blocks))
	ibc.blocks = append(ibc.blocks, NewIBBlock(ibc.ChromosomeName, ibc.BlockSize, blockNum, blockPos, ibc.NumSamples))
	ibc.BlockNames[blockNum] = blockPos
	ibc.NumBlocks++
}

func (ibc *IBChromosome) GetBlock(blockNum uint64) (*IBBlock, bool) {
	if blockPos, ok := ibc.BlockNames[blockNum]; ok {
		if blockPos >= uint64(len(ibc.blocks)) {
			fmt.Println(&ibc, "Index out of range. block num:", blockNum, "block pos:", blockPos, "len:", len(ibc.blocks), "NumBlocks:", ibc.NumBlocks)
			fmt.Println(&ibc, "BlockNames", ibc.BlockNames)
			fmt.Println(&ibc, "Blocks", ibc.blocks)
			debug.PrintStack()
			os.Exit(1)
		}

		return ibc.blocks[blockPos], ok
	} else {
		return &IBBlock{}, ok
	}
}

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) (isNew bool) {
	if _, hasBlock := ibc.GetBlock(blockNum); !hasBlock {
		isNew = false
		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(len(ibc.blocks))

			if NumBlocks != 0 {
				lastBlockPos = NumBlocks - 1
			}

			for currBlockPos := lastBlockPos; currBlockPos < blockNum; currBlockPos++ {
				ibc.AppendBlock(currBlockPos)
			}
		}
		ibc.AppendBlock(blockNum)
	} else {
		isNew = true
	}
	return isNew
}

func (ibc *IBChromosome) Add(reg *interfaces.VCFRegister) (blockNum uint64, isNew bool) {
	position := reg.Position
	distance := reg.Distance
	blockNum = position / ibc.BlockSize

	isNew = ibc.normalizeBlocks(blockNum)

	if block, success := ibc.GetBlock(blockNum); success {
		block.Add(position, distance)
		ibc.Block.Add(position, distance)
		ibc.NumSNPS++
		ibc.MinPosition = tools.Min64(ibc.MinPosition, block.MinPosition)
		ibc.MaxPosition = tools.Max64(ibc.MaxPosition, block.MaxPosition)
	} else {
		fmt.Println("Failure getting block", blockNum)
		os.Exit(1)
	}

	return blockNum, isNew
}

func (ibc *IBChromosome) GenFilename(outPrefix string, format string) (fileName string) {
	baseName := outPrefix + "." + ibc.ChromosomeName
	fileName = save.GenFilename(baseName, format)
	return fileName
}

func (ibc *IBChromosome) Save(outPrefix string, format string) {
	baseName := ibc.GenFilename(outPrefix, format)
	save.Save(baseName, format, ibc)
	ibc.saveBlocks(baseName, format)
}

func (ibc *IBChromosome) saveBlocks(outPrefix string, format string) {
	// BlockNames:     make(map[uint64]uint64, 100),
	// blocks:         make([]*IBBlock, 0, 100),

	for blockNum, blockPos := range ibc.BlockNames {
		block := ibc.blocks[blockPos]

		fmt.Print("saving block: ", outPrefix, " block num: ", blockNum, " block pos: ", blockPos)

		outfile := block.GenFilename(outPrefix, format)

		if _, err := os.Stat(outfile); err == nil {
			// path/to/whatever exists
			fmt.Println(" exists")
			continue

		} else if os.IsNotExist(err) {
			fmt.Println(" creating")
			// path/to/whatever does *not* exist

		} else {
			// Schrodinger: file may or may not exist. See err for details.

			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		}

		block.Save(outPrefix, format)
	}
}
