package ibrowser

import (
	"fmt"
	"math"
	"os"
)

// import "runtime/debug"

//
//
// CHROMOSOME SECTION
//
//

// IBChromosome represents a chromosome
type IBChromosome struct {
	ChromosomeName   string
	ChromosomeNumber int
	BlockSize        uint64
	CounterBits      uint64
	NumSamples       uint64
	MinPosition      uint64
	MaxPosition      uint64
	NumSNPS          uint64
	RegisterSize     uint64
	KeepEmptyBlock   bool
	// BlockNames       map[uint64]uint64
	BlockManager     *BlockManager
	rootBlockManager *BlockManager
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
		" NumBlocks:        ", ibc.NumBlocks(), "\n",
		" NumSNPS:          ", ibc.NumSNPS, "\n",
		" KeepEmptyBlock:   ", ibc.KeepEmptyBlock, "\n",
		// " BlockNames:       ", ibc.BlockNames, "\n",
	)
}

// NewIBChromosome creates a new IBChromosome instance
func NewIBChromosome(
		chromosomeName   string, 
		chromosomeNumber int, 
		blockSize        uint64, 
		counterBits      uint64, 
		numSamples       uint64, 
		keepEmptyBlock   bool,
		rootBlockManager *BlockManager,
	) *IBChromosome {
	
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
		NumSNPS:          0,
		KeepEmptyBlock:   keepEmptyBlock,
		// BlockNames:       make(map[uint64]uint64, 100),
		BlockManager:     NewBlockManager(chromosomeName),
		rootBlockManager: rootBlockManager,
	}

	rootBlockManager.NewBlock(
		chromosomeName,
		chromosomeNumber,
		blockSize,
		counterBits,
		numSamples,
		0,
	)

	return &ibc
}

// AppendBlock appends a block to IBChromosome
func (ibc *IBChromosome) AppendBlock(blockNum uint64) (block *IBBlock) {
	// fmt.Println("IBChromosome :: AppendBlock :: blockNum: ", blockNum)

	if ibc.HasBlock(blockNum) {
		fmt.Println("tried to append existing blockNum:", blockNum)
		os.Exit(1)
	}

	block = 
	ibc.BlockManager.NewBlock(
		ibc.ChromosomeName,
		ibc.ChromosomeNumber,
		ibc.BlockSize,
		ibc.CounterBits,
		ibc.NumSamples,
		blockNum,
	)	

	return block
}

// HasBlock checks whether requested block number exists
func (ibc *IBChromosome) HasBlock(blockNum uint64) bool {
	if _, ok := ibc.BlockManager.GetBlockByNum(blockNum); ok {
		return true
	}
	return false
}

// GetSummaryBlock returns the summar block
func (ibc *IBChromosome) GetSummaryBlock() (block *IBBlock, hasBlock bool) {
	block, hasBlock = ibc.rootBlockManager.GetBlockByName(ibc.ChromosomeName)
	if !hasBlock {
		fmt.Println(ibc.rootBlockManager)
		panic("!GetSummaryBlock")
	}
	return block, hasBlock
}

// GetBlocks returns all blocks
func (ibc *IBChromosome) GetBlocks() ([]*IBBlock, bool) {
	return ibc.BlockManager.Blocks, true
}

// NumBlocks returns the number of blocks
func (ibc *IBChromosome) NumBlocks() (uint64) {
	return ibc.BlockManager.NumBlocks
}

// BlockNumbers returns the block numbers
func (ibc *IBChromosome) BlockNumbers() (map[uint64]uint64) {
	return ibc.BlockManager.BlockNumbers
}

// GetBlock returns one block
func (ibc *IBChromosome) GetBlock(blockNum uint64) (*IBBlock, bool) {
	if block, ok := ibc.BlockManager.GetBlockByNum(blockNum); ok {
		return block, ok
	}

	// debug.PrintStack()

	return nil, false
}

// GetColumn returns the column for a given reference in all blocks
func (ibc *IBChromosome) GetColumn(referenceNumber int) (*[]*IBDistanceTable, bool) {
	cols := make([]*IBDistanceTable, ibc.NumBlocks())
	for bc, block := range ibc.BlockManager.Blocks {
		col, nc := block.GetColumn(referenceNumber)
		if !nc {
			return nil, nc
		}
		cols[bc] = col
	}
	return &cols, false
}

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) (*IBBlock, bool, uint64) {
	block, hasBlock := ibc.GetBlock(blockNum)
	isNew := false
	numBlocksAdded := uint64(0)

	if !hasBlock {
		fmt.Println("  IBChromosome :: normalizeBlocks :: blockNum: ", blockNum, " NEW")

		isNew = true

		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(ibc.NumBlocks())

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

	}

	isNew = false
	return block, isNew, numBlocksAdded
}

// Add adds a new SNP
func (ibc *IBChromosome) Add(reg *VCFRegister) (uint64, bool, uint64) {
	position := reg.Position
	distance := reg.Distance
	blockNum := position / ibc.BlockSize

	block, isNew, numBlocksAdded := ibc.normalizeBlocks(blockNum)

	summaryBlock, hasSummaryBlock := ibc.GetSummaryBlock()
	if !hasSummaryBlock {
		panic("!hasSummaryBlock")
	}

	block.AddVcfMatrix(position, distance)
	summaryBlock.AddVcfMatrix(position, distance)
	ibc.NumSNPS++
	ibc.MinPosition = Min64(ibc.MinPosition, block.MinPosition)
	ibc.MaxPosition = Max64(ibc.MaxPosition, block.MaxPosition)

	return blockNum, isNew, numBlocksAdded
}

//
// Check
//

// Check checks for self consistency
func (ibc *IBChromosome) Check() (res bool) {
	res = true

	res = res && ibc.selfCheck()

	if !res {
		fmt.Printf("Failed chromosome self check\n")
		return res
	}

	for _, block := range ibc.BlockManager.Blocks {
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

	summaryBlock, hasSummaryBlock := ibc.GetSummaryBlock()
	if !hasSummaryBlock {
		panic("!hasSummaryBlock")
	}

	res = res && summaryBlock.Check()

	if !res {
		fmt.Printf("Failed chromosome self check - block chek\n")
		return res
	}

	sumBlock := ibc.GetSumBlocks()

	{
		res = res && (summaryBlock.NumSNPS == sumBlock.NumSNPS)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block NumSNPS: %d != %d\n", ibc.ChromosomeName, summaryBlock.NumSNPS, sumBlock.NumSNPS)
			return res
		}

		res = res && (ibc.NumSNPS == sumBlock.NumSNPS)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumBlock NumSNPS: %d != %d\n", ibc.ChromosomeName, ibc.NumSNPS, sumBlock.NumSNPS)
			return res
		}
	}
	{
		res = res && (summaryBlock.MinPosition == sumBlock.MinPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block MinPosition: %d != %d\n", ibc.ChromosomeName, summaryBlock.MinPosition, sumBlock.MinPosition)
			return res
		}

		res = res && (ibc.MinPosition == sumBlock.MinPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumBlock MinPosition: %d != %d\n", ibc.ChromosomeName, ibc.MinPosition, sumBlock.MinPosition)
			return res
		}

		res = res && (summaryBlock.MaxPosition == sumBlock.MaxPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - block MaxPosition: %d != %d\n", ibc.ChromosomeName, summaryBlock.MaxPosition, sumBlock.MaxPosition)
			return res
		}

		res = res && (ibc.MaxPosition == sumBlock.MaxPosition)

		if !res {
			fmt.Printf("Failed chromosome %s self check - sumblock MaxPosition: %d != %d\n", ibc.ChromosomeName, summaryBlock.MaxPosition, sumBlock.MaxPosition)
			return res
		}
	}

	{
		res = res && (summaryBlock.IsEqual(sumBlock))

		if !res {
			fmt.Printf("Failed chromosome %s self check - blocks not equal\n", ibc.ChromosomeName)
			return res
		}
	}

	return res
}

// GetSumBlocks returns a single block with the sum of all blocks
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

	for _, block := range ibc.BlockManager.Blocks {
		sumBlock.Sum(block)
	}

	return sumBlock
}

//
// Dump
//

// GenMatrixDumpFileName returns the filename of the dump of this project
func (ibc *IBChromosome) GenMatrixDumpFileName(outPrefix string) (filename string) {
	filename = outPrefix + "_chromosomes_" + ibc.ChromosomeName + ".bin"
	return
}

// DumpBlocks dumps all blocks of this chromosome
func (ibc *IBChromosome) DumpBlocks(outPrefix string, isSave bool, isSoft bool) {
	chromosomeFileName := ibc.GenMatrixDumpFileName(outPrefix)

	if isSave {
		ibc.BlockManager.Save(chromosomeFileName)
	} else {
		ibc.BlockManager.Load(chromosomeFileName)
	}
}
