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

	if ibc.HasBlock(blockNum) {
		fmt.Println("tried to append existing blockNum:", blockNum)
		os.Exit(1)
	}

	blockPos := uint64(len(ibc.blocks))

	block = NewIBBlock(
		ibc.ChromosomeName,
		ibc.BlockSize,
		blockNum,
		blockPos,
		ibc.NumSamples)

	ibc.blocks = append(ibc.blocks, block)

	ibc.BlockNames[blockNum] = blockPos

	ibc.NumBlocks = uint64(len(ibc.BlockNames))

	return block
}

func (ibc *IBChromosome) HasBlock(blockNum uint64) bool {
	if _, ok := ibc.BlockNames[blockNum]; ok {
		if ok {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
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

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) (*IBBlock, bool, uint64) {
	// fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum)

	block, hasBlock := ibc.GetBlock(blockNum)
	isNew := false
	numBlocksAdded := uint64(0)

	if !hasBlock {
		fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW")

		isNew = true

		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(len(ibc.blocks))

			if NumBlocks == 0 {
				lastBlockPos = 0
			} else {
				lastBlockPos = NumBlocks
			}

			for currBlockPos := lastBlockPos; currBlockPos < blockNum; currBlockPos++ {
				fmt.Println("IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW. adding intermediate: ", currBlockPos)
				ibc.AppendBlock(currBlockPos)
				numBlocksAdded++
			}
		}

		block := ibc.AppendBlock(blockNum)
		numBlocksAdded++

		return block, isNew, numBlocksAdded

	} else {
		isNew = false
		return block, isNew, numBlocksAdded
	}

	return &IBBlock{}, false, numBlocksAdded
}

func (ibc *IBChromosome) Add(reg *interfaces.VCFRegister) (uint64, bool, uint64) {
	position := reg.Position
	distance := reg.Distance
	blockNum := position / ibc.BlockSize

	block, isNew, numBlocksAdded := ibc.normalizeBlocks(blockNum)

	block.Add(position, distance)
	ibc.block.Add(position, distance)
	ibc.NumSNPS++
	ibc.MinPosition = tools.Min64(ibc.MinPosition, block.MinPosition)
	ibc.MaxPosition = tools.Max64(ibc.MaxPosition, block.MaxPosition)

	return blockNum, isNew, numBlocksAdded
}

func (ibc *IBChromosome) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + ibc.ChromosomeName

	saver := save.NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

func (ibc *IBChromosome) Save(outPrefix string, format string, compression string) {
	baseName, _ := ibc.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)
	saver.Save(ibc)
	ibc.saveBlock(baseName, format, compression)
	ibc.saveBlocks(baseName, format, compression)
}

func (ibc *IBChromosome) Load(outPrefix string, format string, compression string) {
	baseName, _ := ibc.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)
	saver.Load(ibc)
	ibc.loadBlock(baseName, format, compression)
	ibc.loadBlocks(baseName, format, compression)
}

func (ibc *IBChromosome) saveBlock(outPrefix string, format string, compression string) {
	ibc.block.Save(outPrefix+"_block", format, compression)
}

func (ibc *IBChromosome) saveBlocks(outPrefix string, format string, compression string) {
	outPrefix = outPrefix + "_blocks"

	for blockPos := 0; blockPos < len(ibc.blocks); blockPos++ {
		block := ibc.blocks[blockPos]
		blockNum := block.BlockNumber

		fmt.Print("saving block: ", outPrefix, " block num: ", blockNum, " block pos: ", blockPos)

		block.Save(outPrefix, format, compression)
	}
}

func (ibc *IBChromosome) loadBlock(outPrefix string, format string, compression string) {
	ibc.block.Load(outPrefix+"_block", format, compression)
}

func (ibc *IBChromosome) loadBlocks(outPrefix string, format string, compression string) {
	outPrefix = outPrefix + "_blocks"

	for blockPos := 0; blockPos < len(ibc.blocks); blockPos++ {
		block := ibc.blocks[blockPos]
		blockNum := block.BlockNumber

		fmt.Print("loading block: ", outPrefix, " block num: ", blockNum, " block pos: ", blockPos)

		block.Load(outPrefix, format, compression)
	}
}
