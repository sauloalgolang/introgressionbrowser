package ibrowser

import (
	"fmt"
	"math"
	"os"
)

import "runtime/debug"

import (
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	ChromosomeName   string
	ChromosomeNumber int
	BlockSize        uint64
	MinPosition      uint64
	MaxPosition      uint64
	NumBlocks        uint64
	NumSNPS          uint64
	NumSamples       uint64
	CounterBits      int
	KeepEmptyBlock   bool
	BlockNames       map[uint64]uint64
	block            *IBBlock
	blocks           []*IBBlock
}

func NewIBChromosome(chromosomeName string, chromosomeNumber int, blockSize uint64, counterBits int, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	fmt.Println("  NewIBChromosome :: chromosomeName: ", chromosomeName, " chromosomeNumber: ", chromosomeNumber, " blockSize: ", blockSize, " counterBits: ", counterBits, " numSamples: ", numSamples)

	ibc := IBChromosome{
		ChromosomeName:   chromosomeName,
		ChromosomeNumber: chromosomeNumber,
		BlockSize:        blockSize,
		NumSamples:       numSamples,
		MinPosition:      math.MaxUint64,
		MaxPosition:      0,
		NumBlocks:        0,
		NumSNPS:          0,
		CounterBits:      counterBits,
		KeepEmptyBlock:   keepEmptyBlock,
		BlockNames:       make(map[uint64]uint64, 100),
		block:            NewIBBlock("_"+chromosomeName+"_block", blockSize, counterBits, numSamples, 0, 0),
		blocks:           make([]*IBBlock, 0, 100),
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
		ibc.CounterBits,
		ibc.NumSamples,
		blockPos,
		blockNum,
	)

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
		fmt.Println("  IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW")

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

func (ibc *IBChromosome) Add(reg *VCFRegister) (uint64, bool, uint64) {
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

//
// Check
//

func (ibc *IBChromosome) Check() (res bool) {
	res = true

	res = res && ibc.selfCheck()

	if !res {
		return res
	}

	for _, block := range ibc.blocks {
		res = res && block.Check()

		if !res {
			return res
		}
	}

	return res
}

func (ibc *IBChromosome) selfCheck() (res bool) {
	res = true

	res = res && ibc.block.Check()

	if !res {
		return res
	}

	sumBlock := ibc.GetSumBlocks()

	{
		res = res && (ibc.block.NumSNPS == sumBlock.NumSNPS)

		if !res {
			return res
		}

		res = res && (ibc.NumSNPS == sumBlock.NumSNPS)

		if !res {
			return res
		}
	}
	{
		res = res && (ibc.block.MinPosition == sumBlock.MinPosition)

		if !res {
			return res
		}

		res = res && (ibc.MinPosition == sumBlock.MinPosition)

		if !res {
			return res
		}
	}
	{
		res = res && (ibc.block.MaxPosition == sumBlock.MaxPosition)

		if !res {
			return res
		}
		res = res && (ibc.MaxPosition == sumBlock.MaxPosition)

		if !res {
			return res
		}
	}

	{
		res = res && (ibc.block.IsEqual(sumBlock))

		if !res {
			return res
		}
	}

	return res
}

func (ibc *IBChromosome) GetSumBlocks() (sumBlock *IBBlock) {
	sumBlock = NewIBBlock(
		ibc.ChromosomeName,
		ibc.BlockSize,
		ibc.CounterBits,
		ibc.NumSamples,
		0,
		0,
	)

	for _, block := range ibc.blocks {
		sumBlock.Sum(block)
	}

	return sumBlock
}

//
// Filename
//

func (ibc *IBChromosome) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + ibc.ChromosomeName

	saver := save.NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//

func (ibc *IBChromosome) Save(outPrefix string, format string, compression string) {
	ibc.saveLoad(true, outPrefix, format, compression)
}

//
// Load
//
func (ibc *IBChromosome) Load(outPrefix string, format string, compression string) {
	ibc.saveLoad(false, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ibc *IBChromosome) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := ibc.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving chromosome        : ", baseName)
		saver.Save(ibc)
	} else {
		fmt.Println("loading chromosome       : ", baseName)
		saver.Load(ibc)
	}

	ibc.saveLoadBlock(isSave, baseName, format, compression)
	ibc.saveLoadBlocks(isSave, baseName, format, compression)

}

func (ibc *IBChromosome) saveLoadBlock(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_block"

	if isSave {
		fmt.Println("saving chromosome block  : ", newPrefix)
		ibc.block.Save(newPrefix, format, compression)
	} else {
		fmt.Println("loading chromosome block : ", newPrefix)
		ibc.block = NewIBBlock("_"+ibc.ChromosomeName+"_block", ibc.BlockSize, ibc.CounterBits, ibc.NumSamples, 0, 0)
		ibc.block.Load(newPrefix, format, compression)
	}
}

func (ibc *IBChromosome) saveLoadBlocks(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_blocks"
	// fmt.Println("saving blocks", ibc.BlockNames)

	for blockNum, blockPos := range ibc.BlockNames {
		if isSave {
			block := ibc.blocks[blockPos]
			fmt.Printf("saving chromosome blocks :  %-70s block num: %d block pos: %d\n", newPrefix, blockNum, blockPos)
			block.Save(newPrefix, format, compression)

		} else {
			fmt.Printf("loading chromosome blocks:  %-70s block num: %d block pos: %d\n", newPrefix, blockNum, blockPos)

			block := NewIBBlock(
				ibc.ChromosomeName,
				ibc.BlockSize,
				ibc.CounterBits,
				ibc.NumSamples,
				blockPos,
				blockNum,
			)

			ibc.blocks = append(ibc.blocks, block)

			block.Load(newPrefix, format, compression)
		}
	}
}
