package ibrowser

import (
	"fmt"
	"math"
	"os"
)


// BlockManager holds information about blocks in a give domain
type BlockManager struct {
	Blocks       []*IBBlock
	NumBlocks    int64
	domain       string
	BlockNames   map[string]int
	BlockNumbers map[uint64]int
}

// NewBlockManager creates a new BlockManager
func NewBlockManager(domain string) (bm *BlockManager){
	bm = &BlockManager{
		Blocks:       make([]*IBBlock, 0, 0),
		BlockNames:   make(map[string]int, 0),
		BlockNumbers: make(map[uint64]int, 0),
		NumBlocks:    0,
		domain:       domain,
	}

	return bm
}

// NewBlock Registers a new blocks
func (bm *BlockManager) NewBlock(
		chromosomeName   string,
		chromosomeNumber int,
		blockSize        uint64,
		counterBits      uint64,
		numSamples       uint64,
		blockPosition    uint64,
		blockNumber      uint64,
	) (nb *IBBlock) {
	
	bm.BlockNames[chromosomeName] = len(bm.Blocks)
	bm.BlockNumbers[blockNumber ] = len(bm.Blocks)

	nb = NewIBBlock(chromosomeName,
		chromosomeNumber,
		blockSize,
		counterBits,
		numSamples,
		blockPosition,
		blockNumber,
	)

	bm.Blocks = append(bm.Blocks, nb)
	bm.NumBlocks = int64(len(bm.Blocks))

	return nb
}

// GetBlockByName returns block given its name
func (bm *BlockManager) GetBlockByName(name string) (block *IBBlock, ok bool){
	if pos, hasName := bm.BlockNames[name]; hasName {
		block = bm.Blocks[pos]
		return block, true
	}
	return nil, false
}

// GetBlockByNum returns block given its block number
func (bm *BlockManager) GetBlockByNum(blockNum uint64) (block *IBBlock, ok bool){
	if pos, hasNum := bm.BlockNumbers[blockNum]; hasNum {
		block = bm.Blocks[pos]
		return block, true
	}
	return nil, false
}

// GetBlockByPosition returns block given its position
func (bm *BlockManager) GetBlockByPosition(pos uint64) (block *IBBlock, ok bool){
	block = bm.Blocks[pos]
	return block, true
}

// Save to file
func (bm *BlockManager) Save(outPrefix string) {
		// dumperl := NewMultiArrayFile(chromosomeFileName, isSave, isSoft)
	// defer dumperl.Close()

	// ibc.RegisterSize = dumperl.CalculateRegisterSize(ibc.CounterBits, ibc.Block.Matrix.Size)

	// for _, block := range ibc.blockManager.Blocks {
	// 	if isSave {
	// 		block.Dump(dumperl)
	// 	} else {
	// 		block.UnDump(dumperl)
	// 	}
	// }

	// 	dumperg := NewMultiArrayFile(summaryFileName, isSave, isSoft)
	// defer dumperg.Close()

	// ib.RegisterSize = dumperg.CalculateRegisterSize(ib.CounterBits, ib.Block.Matrix.Size)
}

// Load from file
func (bm *BlockManager) Load(outPrefix string) {
}

//
//
// BLOCK SECTION
//
//

// IBBlock holds the information of a block
type IBBlock struct {
	ChromosomeName   string
	ChromosomeNumber int
	BlockSize        uint64
	CounterBits      uint64
	MinPosition      uint64
	MaxPosition      uint64
	NumSNPS          uint64
	NumSamples       uint64
	BlockPosition    uint64
	BlockNumber      uint64
	Serial           uint64
	Matrix           *IBDistanceMatrix
	hasSerial        bool
}

// NewIBBlock generates a new instance of a block
func NewIBBlock(
	chromosomeName string,
	chromosomeNumber int,
	blockSize uint64,
	counterBits uint64,
	numSamples uint64,
	blockPosition uint64,
	blockNumber uint64,
) *IBBlock {
	fmt.Println("   NewIBBlock :: chromosomeName: ", chromosomeName,
		" chromosomeNumber: ", chromosomeNumber,
		" blockSize: ", blockSize,
		" blockPosition: ", blockPosition,
		" blockNumber: ", blockNumber,
		" numSamples: ", numSamples,
	)

	ibb := IBBlock{
		ChromosomeName:   chromosomeName,
		ChromosomeNumber: chromosomeNumber,
		BlockSize:        blockSize,
		CounterBits:      counterBits,
		MinPosition:      math.MaxUint64,
		MaxPosition:      0,
		NumSNPS:          0,
		NumSamples:       numSamples,
		BlockPosition:    blockPosition,
		BlockNumber:      blockNumber,
		Serial:           0,
		hasSerial:        false,
		Matrix: NewDistanceMatrix(
			chromosomeName,
			blockSize,
			counterBits,
			numSamples,
			blockPosition,
			blockNumber,
		),
	}

	return &ibb
}

func (ibb *IBBlock) String() string {
	return fmt.Sprint("Block :: ",
		" ChromosomeName:   ", ibb.ChromosomeName, "\n",
		" ChromosomeNumber: ", ibb.ChromosomeNumber, "\n",
		" BlockSize:        ", ibb.BlockSize, "\n",
		" CounterBits:      ", ibb.CounterBits, "\n",
		" MinPosition:      ", ibb.MinPosition, "\n",
		" MaxPosition:      ", ibb.MaxPosition, "\n",
		" NumSNPS:          ", ibb.NumSNPS, "\n",
		" NumSamples:       ", ibb.NumSamples, "\n",
		" BlockPosition:    ", ibb.BlockPosition, "\n",
		" BlockNumber:      ", ibb.BlockNumber, "\n",
		" Serial:           ", ibb.Serial, "\n",
	)
}

// AddVcfMatrix add a vcf matrix to the blocks
func (ibb *IBBlock) AddVcfMatrix(position uint64, distance *VCFDistanceMatrix) {
	// fmt.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = Min64(ibb.MinPosition, position)
	ibb.MaxPosition = Max64(ibb.MaxPosition, position)
	ibb.Matrix.IncrementWithVcfMatrix(distance)
}

// Add add another ibmatrix to the blocks
func (ibb *IBBlock) Add(position uint64, distance *IBDistanceMatrix) {
	// fmt.Println("Add", position, ibb.NumSNPS, ibb)
	ibb.NumSNPS++
	ibb.MinPosition = Min64(ibb.MinPosition, position)
	ibb.MaxPosition = Max64(ibb.MaxPosition, position)
	ibb.Matrix.Merge(distance)
}

// GetMatrix gets the summary matrix
func (ibb *IBBlock) GetMatrix() (*IBDistanceMatrix, bool) {
	return ibb.Matrix, true
}

// GetMatrixTable gets the table from the summary matrix
func (ibb *IBBlock) GetMatrixTable() (*IBDistanceTable, bool) {
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		return nil, false
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, false
	}

	return table, true
}

// GetColumn returns a give column from the table
func (ibb *IBBlock) GetColumn(referenceNumber int) (*IBDistanceTable, bool) {
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		return nil, false
	}

	col, hasCol := matrix.GetColumn(referenceNumber)

	if !hasCol {
		return nil, false
	}

	return col, hasCol
}

// Sum returns the sum between another block and this block
func (ibb *IBBlock) Sum(other *IBBlock) {
	ibb.NumSNPS += other.NumSNPS
	ibb.MinPosition = Min64(ibb.MinPosition, other.MinPosition)
	ibb.MaxPosition = Max64(ibb.MaxPosition, other.MaxPosition)

	matrix, hasMatrix := other.GetMatrix()

	if !hasMatrix {
		panic("no matrix")
	}

	ibb.Matrix.Merge(matrix)
}

// IsEqual returns wether this block is equal to the other block
func (ibb *IBBlock) IsEqual(other *IBBlock) (res bool) {
	res = true

	res = res && (ibb.NumSNPS == other.NumSNPS)

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - NumSNPS: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.NumSNPS, other.NumSNPS)
		return res
	}

	if ibb.NumSNPS > 0 {
		res = res && (ibb.MinPosition == other.MinPosition)

		if !res {
			fmt.Printf("IsEqual :: Failed block %s - #%d check - MinPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, other.MinPosition)
			return res
		}

		res = res && (ibb.MaxPosition == other.MaxPosition)

		if !res {
			fmt.Printf("IsEqual :: Failed block %s - #%d check - MaxPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MaxPosition, other.MaxPosition)
			return res
		}
	}

	matrix, _ := other.GetMatrix()

	res = res && ibb.Matrix.IsEqual(matrix)

	if !res {
		fmt.Printf("IsEqual :: Failed block %s - #%d check - Matrix not equal\n", ibb.ChromosomeName, ibb.BlockNumber)
		return res
	}

	return res
}

//
// Serial
//

// SetSerial sets the serial number of this block
func (ibb *IBBlock) SetSerial(serial uint64) {
	ibb.Serial = serial
	ibb.Matrix.Serial = serial
}

// CheckSerial checks wether the serial numbers match
func (ibb *IBBlock) CheckSerial(serial uint64) bool {
	eq1 := ibb.Serial == serial

	if !eq1 {
		fmt.Println("block serial ", ibb.Serial, " != ", serial, ibb)
	} else {
		// fmt.Println("block serial ", ibb.Serial, " == ", serial, ibb)
	}

	eq2 := ibb.Matrix.Serial == serial

	if !eq2 {
		fmt.Println("matrix serial ", ibb.Matrix.Serial, " != ", serial, ibb.Matrix)
	} else {
		// fmt.Println("matrix serial ", ibb.Matrix.Serial, " == ", serial, ibb.Matrix)
	}

	return eq1 && eq2
}

//
// Check
//

// Check checks the self consistency of the data
func (ibb *IBBlock) Check() (res bool) {
	res = true

	res = res && (ibb.BlockNumber == ibb.Matrix.BlockNumber)

	if !res {
		fmt.Printf("Failed block %s - #%d check - BlockNumber: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockNumber, ibb.Matrix.BlockNumber)
		return res
	}

	res = res && (ibb.BlockPosition == ibb.Matrix.BlockPosition)

	if !res {
		fmt.Printf("Failed block %s - #%d check - BlockPosition: %d != %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.BlockPosition, ibb.Matrix.BlockPosition)
		return res
	}

	// res = res && (ibb.ChromosomeName == ibb.Matrix.ChromosomeName)

	// if !res {
	// 	fmt.Printf("Failed block %s - #%d check - ChromosomeName: %s != %s\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.ChromosomeName, ibb.Matrix.ChromosomeName)
	// 	return res
	// }

	if ibb.NumSNPS > 0 {
		res = res && (ibb.MinPosition <= ibb.MaxPosition)

		if !res {
			fmt.Printf("Failed block %s - #%d check - MinPosition %d > MaxPosition %d\n", ibb.ChromosomeName, ibb.BlockNumber, ibb.MinPosition, ibb.MaxPosition)
			fmt.Println(ibb)
			return res
		}
	}

	return res
}

//
// Filename
//

// GenFilename generates the filename to save this block
func (ibb *IBBlock) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix + "." + fmt.Sprintf("%012d", ibb.BlockNumber)

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//

// Save saves this block to file
func (ibb *IBBlock) Save(outPrefix string, format string, compression string) {
	ibb.saveLoad(true, outPrefix, format, compression)
}

//
// Load
//

// Load loads this command from file
func (ibb *IBBlock) Load(outPrefix string, format string, compression string) {
	ibb.saveLoad(false, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ibb *IBBlock) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := ibb.GenFilename(outPrefix, format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Printf("saving block             :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Save(ibb)
		ibb.Matrix.Save(baseName, format, compression)
	} else {
		fmt.Printf("loading block            :  %-70s block num: %d block pos: %d\n", baseName, ibb.BlockNumber, ibb.BlockPosition)
		saver.Load(ibb)

		ibb.Matrix = NewDistanceMatrix(
			ibb.ChromosomeName,
			ibb.BlockSize,
			ibb.CounterBits,
			ibb.NumSamples,
			ibb.BlockPosition,
			ibb.BlockNumber,
		)

		ibb.Matrix.Load(baseName, format, compression)
	}
}

//
// Dump
//

// Dump dumps the matrix table to a binary file
func (ibb *IBBlock) Dump(dumper *MultiArrayFile) (serial uint64) {
	serial = uint64(0)
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		fmt.Println("failed getting matrix")
		os.Exit(1)
	}

	serial = matrix.Dump(dumper)
	ibb.SetSerial(serial)

	return
}

// UnDump reads the matrix table from a binary file
func (ibb *IBBlock) UnDump(dumper *MultiArrayFile) (serial uint64, hasData bool) {
	matrix, hasMatrix := ibb.GetMatrix()

	if !hasMatrix {
		fmt.Println("failed getting matrix")
		os.Exit(1)
	}

	serial, hasData = matrix.UnDump(dumper)

	if !hasData {
		fmt.Println("Tried to read beyond the file")
		os.Exit(1)
	}

	if !ibb.CheckSerial(serial) {
		fmt.Println("Mismatch in order of files")
		os.Exit(1)
	}

	return
}
