package ibrowser

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
	if pos, hasName := bm.BlockNames[name]; hasName {
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

	for _, block := range bm.Blocks {
		if isSave {
			block.Dump(dumper)
		} else {
			block.UnDump(dumper)
		}
	}
}
