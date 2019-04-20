package ibrowser

import (
	"fmt"
	"math"
	"os"
)

import "runtime/debug"

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	ChromosomeName   string
	ChromosomeNumber int
	BlockSize        uint64
	CounterBits      int
	NumSamples       uint64
	MinPosition      uint64
	MaxPosition      uint64
	NumBlocks        uint64
	NumSNPS          uint64
	KeepEmptyBlock   bool
	BlockNames       map[uint64]uint64
	Block            *IBBlock
	Blocks           []*IBBlock
}

func (ibc *IBChromosome) String() string {
	return fmt.Sprint("Block :: ",
		" ChromosomeName:   ", ibc.ChromosomeName, "\n",
		" ChromosomeNumber: ", ibc.ChromosomeNumber, "\n",
		" BlockSize:        ", ibc.BlockSize, "\n",
		" CounterBits:      ", ibc.CounterBits, "\n",
		" NumSamples:       ", ibc.NumSamples, "\n",
		" MinPosition:      ", ibc.MinPosition, "\n",
		" MaxPosition:      ", ibc.MaxPosition, "\n",
		" NumBlocks:        ", ibc.NumBlocks, "\n",
		" NumSNPS:          ", ibc.NumSNPS, "\n",
		" KeepEmptyBlock:   ", ibc.KeepEmptyBlock, "\n",
		" NumSBlockNames:   ", ibc.BlockNames, "\n",
	)
}

func NewIBChromosome(chromosomeName string, chromosomeNumber int, blockSize uint64, counterBits int, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	fmt.Println("  NewIBChromosome :: chromosomeName: ", chromosomeName,
		" chromosomeNumber: ", chromosomeNumber,
		" blockSize: ", blockSize,
		" counterBits: ", counterBits,
		" numSamples: ", numSamples,
	)

	ibc := IBChromosome{
		ChromosomeName:   chromosomeName,
		ChromosomeNumber: chromosomeNumber,
		BlockSize:        blockSize,
		CounterBits:      counterBits,
		NumSamples:       numSamples,
		MinPosition:      math.MaxUint64,
		MaxPosition:      0,
		NumBlocks:        0,
		NumSNPS:          0,
		KeepEmptyBlock:   keepEmptyBlock,
		BlockNames:       make(map[uint64]uint64, 100),
		Block:            NewIBBlock("_"+chromosomeName+"_block", chromosomeNumber, blockSize, counterBits, numSamples, 0, 0),
		Blocks:           make([]*IBBlock, 0, 100),
	}

	return &ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) (block *IBBlock) {
	// fmt.Println("IBChromosome :: AppendBlock :: blockNum: ", blockNum)

	if ibc.HasBlock(blockNum) {
		fmt.Println("tried to append existing blockNum:", blockNum)
		os.Exit(1)
	}

	blockPos := uint64(len(ibc.Blocks))

	block = NewIBBlock(
		ibc.ChromosomeName,
		ibc.ChromosomeNumber,
		ibc.BlockSize,
		ibc.CounterBits,
		ibc.NumSamples,
		blockPos,
		blockNum,
	)

	ibc.Blocks = append(ibc.Blocks, block)

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
			NumBlocks := uint64(len(ibc.Blocks))

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
	ibc.Block.Add(position, distance)
	ibc.NumSNPS++
	ibc.MinPosition = Min64(ibc.MinPosition, block.MinPosition)
	ibc.MaxPosition = Max64(ibc.MaxPosition, block.MaxPosition)

	return blockNum, isNew, numBlocksAdded
}

//
// Check
//

func (ibc *IBChromosome) Check() (res bool) {
	res = true

	res = res && ibc.selfCheck()

	if !res {
		fmt.Printf("Failed chromosome self check\n")
		return res
	}

	for _, block := range ibc.Blocks {
		res = res && block.Check()

		if !res {
			fmt.Printf("Failed chromosome - block check - block %s pos %d number %d\n",
				block.ChromosomeName,
				block.BlockPosition,
				block.BlockNumber,
			)
			return res
		}
	}

	return res
}

func (ibc *IBChromosome) selfCheck() (res bool) {
	res = true

	res = res && ibc.Block.Check()

	if !res {
		fmt.Printf("Failed chromosome self check - block chek\n")
		return res
	}

	sumBlock := ibc.GetSumBlocks()

	{
		res = res && (ibc.Block.NumSNPS == sumBlock.NumSNPS)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block NumSNPS: %d != %d\n", ibc.ChromosomeName, ibc.Block.NumSNPS, sumBlock.NumSNPS)
			return res
		}

		res = res && (ibc.NumSNPS == sumBlock.NumSNPS)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumBlock NumSNPS: %d != %d\n", ibc.ChromosomeName, ibc.NumSNPS, sumBlock.NumSNPS)
			return res
		}
	}
	{
		res = res && (ibc.Block.MinPosition == sumBlock.MinPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block MinPosition: %d != %d\n", ibc.ChromosomeName, ibc.Block.MinPosition, sumBlock.MinPosition)
			return res
		}

		res = res && (ibc.MinPosition == sumBlock.MinPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumBlock MinPosition: %d != %d\n", ibc.ChromosomeName, ibc.MinPosition, sumBlock.MinPosition)
			return res
		}

		res = res && (ibc.Block.MaxPosition == sumBlock.MaxPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block MaxPosition: %d != %d\n", ibc.ChromosomeName, ibc.Block.MaxPosition, sumBlock.MaxPosition)
			return res
		}

		res = res && (ibc.MaxPosition == sumBlock.MaxPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumblock MaxPosition: %d != %d\n", ibc.ChromosomeName, ibc.Block.MaxPosition, sumBlock.MaxPosition)
			return res
		}
	}

	{
		res = res && (ibc.Block.IsEqual(sumBlock))

		if !res {
			fmt.Printf("Failed chromosome %s self check - blocks not equal\n", ibc.ChromosomeName)
			return res
		}
	}

	return res
}

func (ibc *IBChromosome) GetSumBlocks() (sumBlock *IBBlock) {
	sumBlock = NewIBBlock(
		ibc.ChromosomeName,
		ibc.ChromosomeNumber,
		ibc.BlockSize,
		ibc.CounterBits,
		ibc.NumSamples,
		0,
		0,
	)

	for _, block := range ibc.Blocks {
		sumBlock.Sum(block)
	}

	return sumBlock
}

//
// Filename
//

func (ibc *IBChromosome) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + ibc.ChromosomeName

	saver := NewSaverCompressed(baseName, format, compression)

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
	saver := NewSaverCompressed(baseName, format, compression)

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
		ibc.Block.Save(newPrefix, format, compression)
	} else {
		fmt.Println("loading chromosome block : ", newPrefix)
		ibc.Block = NewIBBlock(
			"_"+ibc.ChromosomeName+"_block",
			ibc.ChromosomeNumber,
			ibc.BlockSize,
			ibc.CounterBits,
			ibc.NumSamples,
			0,
			0,
		)
		ibc.Block.Load(newPrefix, format, compression)
	}
}

func (ibc *IBChromosome) saveLoadBlocks(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_blocks"
	// fmt.Println("saving blocks", ibc.BlockNames)

	for blockPos, block := range ibc.Blocks {
		blockNum := block.BlockNumber

		if isSave {
			block := ibc.Blocks[blockPos]
			fmt.Printf("saving chromosome blocks :  %-70s block num: %d block pos: %d\n", newPrefix, blockNum, blockPos)
			block.Save(newPrefix, format, compression)

		} else {
			fmt.Printf("loading chromosome blocks:  %-70s block num: %d block pos: %d\n", newPrefix, blockNum, blockPos)

			block := NewIBBlock(
				ibc.ChromosomeName,
				ibc.ChromosomeNumber,
				ibc.BlockSize,
				ibc.CounterBits,
				ibc.NumSamples,
				uint64(blockPos),
				blockNum,
			)

			ibc.Blocks = append(ibc.Blocks, block)

			block.Load(newPrefix, format, compression)
		}
	}
}
