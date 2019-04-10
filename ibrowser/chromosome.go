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
	BlockNames     map[uint64]uint64
	block          *IBBlock
	blocks         []*IBBlock
}

func NewIBChromosome(chromosomeName string, blockSize uint64, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	fmt.Println("NewIBChromosome :: chromosomeName: ", chromosomeName, " blockSize: ", blockSize, " numSamples: ", numSamples)

	ibc := IBChromosome{
		ChromosomeName: chromosomeName,
		BlockSize:      blockSize,
		NumSamples:     numSamples,
		MinPosition:    math.MaxUint64,
		MaxPosition:    0,
		NumBlocks:      0,
		NumSNPS:        0,
		KeepEmptyBlock: keepEmptyBlock,
		BlockNames:     make(map[uint64]uint64, 100),
		block:          NewIBBlock("_"+chromosomeName+"_block", blockSize, 0, 0, numSamples),
		blocks:         make([]*IBBlock, 0, 100),
	}

	return &ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) (block *IBBlock) {
	// fmt.Println("IBChromosome :: AppendBlock :: blockNum: ", blockNum)

	blockPos := uint64(len(ibc.blocks))

	block = NewIBBlock(
		ibc.ChromosomeName,
		ibc.BlockSize,
		blockNum,
		blockPos,
		ibc.NumSamples)

	ibc.blocks = append(ibc.blocks, block)

	ibc.BlockNames[blockNum] = blockPos
	ibc.NumBlocks++

	return block
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

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) (*IBBlock, bool) {
	// fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum)

	block, hasBlock := ibc.GetBlock(blockNum)
	isNew := false

	if !hasBlock {
		fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW")

		isNew = true

		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(len(ibc.blocks))

			if NumBlocks == 0 {
				lastBlockPos = 0
			} else {
				lastBlockPos = NumBlocks - 1
			}

			for currBlockPos := lastBlockPos; currBlockPos < blockNum; currBlockPos++ {
				fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW. adding intermediate")
				ibc.AppendBlock(currBlockPos)
			}
		}

		block := ibc.AppendBlock(blockNum)

		return block, isNew

	} else {
		isNew = false
		return block, isNew
	}

	return &IBBlock{}, false
}

func (ibc *IBChromosome) Add(reg *interfaces.VCFRegister) (uint64, bool) {
	position := reg.Position
	distance := reg.Distance
	blockNum := position / ibc.BlockSize

	block, isNew := ibc.normalizeBlocks(blockNum)

	block.Add(position, distance)
	ibc.block.Add(position, distance)
	ibc.NumSNPS++
	ibc.MinPosition = tools.Min64(ibc.MinPosition, block.MinPosition)
	ibc.MaxPosition = tools.Max64(ibc.MaxPosition, block.MaxPosition)

	return blockNum, isNew
}

func (ibc *IBChromosome) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix + "." + ibc.ChromosomeName
	fileName = save.GenFilename(baseName, format)
	return baseName, fileName
}

func (ibc *IBChromosome) Save(outPrefix string, format string) {
	baseName, _ := ibc.GenFilename(outPrefix, format)
	save.Save(baseName, format, ibc)
	ibc.saveBlock(baseName, format)
	ibc.saveBlocks(baseName, format)
}

func (ibc *IBChromosome) saveBlock(outPrefix string, format string) {
	ibc.block.Save(outPrefix+"_block", format)
}

func (ibc *IBChromosome) saveBlocks(outPrefix string, format string) {
	// BlockNames:     make(map[uint64]uint64, 100),
	// blocks:         make([]*IBBlock, 0, 100),

	outPrefix = outPrefix + "_blocks"

	for blockNum, blockPos := range ibc.BlockNames {
		block := ibc.blocks[blockPos]

		_, fileName := block.GenFilename(outPrefix, format)

		fmt.Print("saving block: ", outPrefix, " block num: ", blockNum, " block pos: ", blockPos, " to: ", fileName)

		if _, err := os.Stat(fileName); err == nil {
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
