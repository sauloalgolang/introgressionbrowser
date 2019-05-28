package ibrowser

import (
	"fmt"
	// log "github.com/sirupsen/logrus"
)

// BlockManager holds information about blocks in a give domain
type BlockManager struct {
	Domain       string
	FileName     string
	CounterBits  uint64
	NumSamples   uint64
	NumBlocks    uint64
	Dimension    uint64
	BlockSize    uint64
	Blocks       []*IBBlock
	BlockNames   map[string]uint64
	BlockNumbers map[uint64]uint64
	Bmm          BlockManagerMMap
}

// NewBlockManager creates a new BlockManager
func NewBlockManager(
	domain string,
	fileName string,
	counterBits uint64,
	numSamples uint64,
	blockSize uint64,
) (bm *BlockManager) {

	bm = &BlockManager{
		Domain:       domain,
		FileName:     fileName,
		CounterBits:  counterBits,
		NumSamples:   numSamples,
		BlockSize:    blockSize,
		NumBlocks:    0,
		Dimension:    0,
		Blocks:       make([]*IBBlock, 0, 0),
		BlockNames:   make(map[string]uint64, 0),
		BlockNumbers: make(map[uint64]uint64, 0),
	}

	tm := NewTriangularMatrix(numSamples)
	dimension := tm.CalculateSize()
	bmm, err := NewBlockManagerMMap(fileName, counterBits, dimension, RW)

	if err != nil {
		return nil
	}

	bm.Bmm = *bmm
	bm.Dimension = dimension

	return bm
}

func (bm *BlockManager) String() (res string) {
	res += fmt.Sprintf("BlockManager ::\n")
	res += fmt.Sprintf(" Domain       %s\n", bm.Domain)
	res += fmt.Sprintf(" File Name    %d\n", bm.FileName)
	res += fmt.Sprintf(" CounterBits  %d\n", bm.CounterBits)
	res += fmt.Sprintf(" Num Samples  %d\n", bm.NumSamples)
	res += fmt.Sprintf(" Num Blocks   %d\n", bm.NumBlocks)
	res += fmt.Sprintf(" Dimension    %d\n", bm.Dimension)
	res += fmt.Sprintf(" BlockSize    %d\n", bm.BlockSize)
	res += fmt.Sprintf(" Blocks #     %d\n", len(bm.Blocks))
	res += fmt.Sprintf(" BlockNames   %#v\n", bm.BlockNames)
	res += fmt.Sprintf(" BlockNumbers %#v\n", bm.BlockNumbers)
	res += fmt.Sprintf(" MMaper       %#v\n", bm.Bmm)

	return
}

// NewBlock Registers a new blocks
func (bm *BlockManager) NewBlock(
	chromosomeName string,
	chromosomeNumber int,
	blockNumber uint64,
) (nb *IBBlock) {
	blockPosition := uint64(len(bm.Blocks))

	bm.BlockNames[chromosomeName] = blockPosition
	bm.BlockNumbers[blockNumber] = blockPosition
	gn := bm.GetMatrixMaker()

	nb = NewIBBlock(
		chromosomeName,
		chromosomeNumber,
		blockNumber,
		blockPosition,
		bm.BlockSize,
		bm.CounterBits,
		bm.NumSamples,
		gn,
	)

	bm.Blocks = append(bm.Blocks, nb)
	bm.NumBlocks = uint64(len(bm.Blocks))

	return nb
}

// GetMatrixMaker - Return default matrix maker
func (bm *BlockManager) GetMatrixMaker() MatrixMaker {
	return bm.Bmm.GetMatrixMaker()
}

// GetFallbackMatrixMaker - Return fallback matrix maker
func (bm *BlockManager) GetFallbackMatrixMaker() MatrixMaker {
	return bm.Bmm.GetFallbackMatrixMaker()
}

// GetBlockByName returns block given its name
func (bm *BlockManager) GetBlockByName(name string) (block *IBBlock, ok bool) {
	// log.Println("BlockManager.GetBlockByName :: name:", name)
	// log.Println("BlockManager.GetBlockByName :: bm:", bm)
	// log.Println("BlockManager.GetBlockByName :: BlockNames:", bm.BlockNames)
	// log.Println("BlockManager.GetBlockByName :: Blocks", bm.Blocks)
	if pos, hasName := bm.BlockNames[name]; hasName {
		if pos > uint64(len(bm.Blocks)) {
			panic(fmt.Sprintf("BlockManager.GetBlockByName :: access error. position %d > %d", pos, len(bm.Blocks)))
		}
		block = bm.Blocks[pos]
		return block, true
	}
	return nil, false
}

// GetBlockByNum returns block given its block number
func (bm *BlockManager) GetBlockByNum(blockNum uint64) (block *IBBlock, ok bool) {
	if pos, hasNum := bm.BlockNumbers[blockNum]; hasNum {
		block = bm.Blocks[pos]
		return block, true
	}
	return nil, false
}

// GetBlockByPosition returns block given its position
func (bm *BlockManager) GetBlockByPosition(pos uint64) (block *IBBlock, ok bool) {
	block = bm.Blocks[pos]
	return block, true
}

// // Save to file
// func (bm *BlockManager) Save(outPrefix string) {
// 	bm.saveLoad(outPrefix, true)
// }

// // Load from file
// func (bm *BlockManager) Load(outPrefix string) {
// 	bm.saveLoad(outPrefix, false)
// }

// func (bm *BlockManager) saveLoad(outPrefix string, isSave bool) {
// 	isSoft := false
// 	dumper := NewMultiArrayFile(outPrefix, isSave, isSoft)
// 	defer dumper.Close()

// 	// log.Println(bm)

// 	for b, block := range bm.Blocks {
// 		log.Println("BlockManager.saveLoad", b)
// 		// log.Println("BlockManager.saveLoad", b, " - BM - ", bm)
// 		if isSave {
// 			block.Dump(dumper)
// 		} else {
// 			block.UnDump(dumper)
// 		}
// 		// log.Println("BlockManager.saveLoad", b, " - DONE")
// 		// log.Println("BlockManager.saveLoad", b, " - BLOCK", bm.Blocks)
// 	}
// }
