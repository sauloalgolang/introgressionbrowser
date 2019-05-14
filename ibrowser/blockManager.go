package ibrowser

import "fmt"

// BlockManager holds information about blocks in a give domain
type BlockManager struct {
	Domain       string
	Blocks       []*IBBlock
	NumBlocks    uint64
	BlockNames   map[string]uint64
	BlockNumbers map[uint64]uint64
}

// NewBlockManager creates a new BlockManager
func NewBlockManager(domain string) (bm *BlockManager) {
	bm = &BlockManager{
		Domain:       domain,
		Blocks:       make([]*IBBlock, 0, 0),
		BlockNames:   make(map[string]uint64, 0),
		BlockNumbers: make(map[uint64]uint64, 0),
		NumBlocks:    0,
	}

	return bm
}

func (bm *BlockManager) String() (res string) {
	res += fmt.Sprintf("BlockManager ::\n")
	res += fmt.Sprintf(" Domain       %s\n", bm.Domain)
	res += fmt.Sprintf(" Blocks #     %d\n", len(bm.Blocks))
	res += fmt.Sprintf(" Num Blocks   %d\n", bm.NumBlocks)
	res += fmt.Sprintf(" BlockNames   %#v\n", bm.BlockNames)
	res += fmt.Sprintf(" BlockNumbers %#v\n", bm.BlockNumbers)
	return
}

// NewBlock Registers a new blocks
func (bm *BlockManager) NewBlock(
	chromosomeName string,
	chromosomeNumber int,
	blockSize uint64,
	counterBits uint64,
	numSamples uint64,
	blockNumber uint64,
) (nb *IBBlock) {

	blockPosition := uint64(len(bm.Blocks))
	bm.BlockNames[chromosomeName] = blockPosition
	bm.BlockNumbers[blockNumber] = blockPosition

	nb = NewIBBlock(chromosomeName,
		chromosomeNumber,
		blockSize,
		counterBits,
		numSamples,
		blockNumber,
		blockPosition,
	)

	bm.Blocks = append(bm.Blocks, nb)
	bm.NumBlocks = uint64(len(bm.Blocks))

	return nb
}

// GetBlockByName returns block given its name
func (bm *BlockManager) GetBlockByName(name string) (block *IBBlock, ok bool) {
	// fmt.Println("BlockManager.GetBlockByName :: name:", name)
	// fmt.Println("BlockManager.GetBlockByName :: bm:", bm)
	// fmt.Println("BlockManager.GetBlockByName :: BlockNames:", bm.BlockNames)
	// fmt.Println("BlockManager.GetBlockByName :: Blocks", bm.Blocks)
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

// Save to file
func (bm *BlockManager) Save(outPrefix string) {
	bm.saveLoad(outPrefix, true)
}

// Load from file
func (bm *BlockManager) Load(outPrefix string) {
	bm.saveLoad(outPrefix, false)
}

func (bm *BlockManager) saveLoad(outPrefix string, isSave bool) {
	isSoft := false
	dumper := NewMultiArrayFile(outPrefix, isSave, isSoft)
	defer dumper.Close()

	// fmt.Println(bm)

	for b, block := range bm.Blocks {
		fmt.Println("BlockManager.saveLoad", b)
		// fmt.Println("BlockManager.saveLoad", b, " - BM - ", bm)
		if isSave {
			block.Dump(dumper)
		} else {
			block.UnDump(dumper)
		}
		// fmt.Println("BlockManager.saveLoad", b, " - DONE")
		// fmt.Println("BlockManager.saveLoad", b, " - BLOCK", bm.Blocks)
	}
}
