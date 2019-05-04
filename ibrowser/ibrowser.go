package ibrowser

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/save"
)

var mutex = &sync.Mutex{}

//
//
// IBROWSER SECTION
//
//

// IBrowser representa the whole ibrowser database
type IBrowser struct {
	Samples        VCFSamples
	NumSamples     uint64
	BlockSize      uint64
	CounterBits    uint64
	KeepEmptyBlock bool
	NumRegisters   uint64
	NumSNPS        uint64
	NumBlocks      uint64
	RegisterSize   uint64
	Parameters     Parameters
	//
	ChromosomesNames NamePosPairList
	Chromosomes      map[string]*IBChromosome
	Block            *IBBlock
	//
	lastChrom    string
	lastPosition uint64
	blockManager *BlockManager
	//
	// Header string
	//
	// TODO: per sample stats
}

// NewIBrowser generates a new instance of IBrowser
func NewIBrowser(parameters Parameters) *IBrowser {
	blockSize := parameters.BlockSize
	counterBits := parameters.CounterBits
	keepEmptyBlock := parameters.KeepEmptyBlock

	if blockSize > uint64((math.MaxUint32/3)-1) {
		fmt.Println("block size too large")
		os.Exit(1)
	}

	ib := IBrowser{
		Samples:    make(VCFSamples, 0, 100),
		NumSamples: 0,
		//
		BlockSize:      blockSize,
		KeepEmptyBlock: keepEmptyBlock,
		//
		NumRegisters: 0,
		NumSNPS:      0,
		NumBlocks:    0,
		CounterBits:  counterBits,
		RegisterSize: 0,
		Parameters:   parameters,
		//
		lastChrom:    "",
		lastPosition: 0,
		//
		ChromosomesNames: make(NamePosPairList, 0, 100),
		Chromosomes:      make(map[string]*IBChromosome, 100),
		//
		// block: NewIBBlock("_whole_genome", blockSize, counterBits, 0, 0, 0),
		blockManager: NewBlockManager("_whole_genome"),
	}

	return &ib
}

// SetSamples set the sample names
func (ib *IBrowser) SetSamples(samples *VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(VCFSamples, numSamples, numSamples)

	ib.NumSamples = uint64(numSamples)
	ib.Block = ib.blockManager.NewBlock("_whole_genome", 0, ib.BlockSize, ib.CounterBits, ib.NumSamples, 0, 0)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

// GetOrCreateChromosome gets or creates a new chromosome if it does not already exists
func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "exists", &chromosome)
		return chromosome
	}
	// fmt.Println("GetOrCreateChromosome", chromosomeName, "creating")
	return ib.AddChromosome(chromosomeName, chromosomeNumber)
}

// AddChromosome adds a new chromosome
func (ib *IBrowser) AddChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, hasChromosome := ib.GetChromosome(chromosomeName); hasChromosome {
		fmt.Println("Failed to add chromosome", chromosomeName, ". Already exists", &chromosome)
		os.Exit(1)
	}

	ib.Chromosomes[chromosomeName] = NewIBChromosome(
		chromosomeName,
		chromosomeNumber,
		ib.BlockSize,
		ib.CounterBits,
		ib.NumSamples,
		ib.KeepEmptyBlock,
		ib.blockManager,
	)

	ib.ChromosomesNames = append(ib.ChromosomesNames, NamePosPair{chromosomeName, chromosomeNumber})

	sort.Sort(ib.ChromosomesNames)

	return ib.Chromosomes[chromosomeName]
}

// RegisterCallBack is the callback function to receive a VCF register
func (ib *IBrowser) RegisterCallBack(samples *VCFSamples, reg *VCFRegister) {
	if atomic.LoadUint64(&ib.NumSamples) == 0 {
		ib.SetSamples(samples)

	} else {
		if len(ib.Samples) != len(*samples) {
			fmt.Println("Sample mismatch")
			fmt.Println(len(ib.Samples), "!=", len(*samples))
			os.Exit(1)
		}
	}

	chromosome := ib.GetOrCreateChromosome(reg.Chromosome, reg.ChromosomeNumber)

	_, isNew, numBlocksAdded := chromosome.Add(reg)

	mutex.Lock()
	{
		if isNew {
			ib.NumBlocks += numBlocksAdded
		}

		ib.NumRegisters++

		ib.NumSNPS++

		ib.Block.AddVcfMatrix(0, reg.Distance)
	}
	mutex.Unlock()
}

// Check checks for self consistency in the data
func (ib *IBrowser) Check() (res bool) {
	fmt.Println("Starting self check")

	res = true

	res = res && ib.selfCheck()

	if !res {
		fmt.Println("Failed ibrowser self check")
		return res
	}

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.Chromosomes[chromosomeName.Name]

		res = res && chromosome.Check()

		if !res {
			fmt.Printf("Failed ibrowser chromosome %s (%d) check\n", chromosomeName.Name, chromosomeName.Pos)
			return res
		}
	}

	return res
}

func (ib *IBrowser) selfCheck() (res bool) {
	res = true

	res = res && ib.Block.Check()

	if !res {
		fmt.Printf("Failed ibrowser self check - block check\n")
		return res
	}

	{
		res = res && (ib.BlockSize == ib.Block.BlockSize)

		if !res {
			fmt.Printf("Failed ibrowser self check - block size %d != %d\n", ib.BlockSize, ib.Block.BlockSize)
			return res
		}

		res = res && (ib.CounterBits == ib.Block.CounterBits)

		if !res {
			fmt.Printf("Failed ibrowser self check - CounterBits %d != %d\n", ib.CounterBits, ib.Block.CounterBits)
			return res
		}
	}

	res = res && (ib.NumSNPS == ib.Block.NumSNPS)

	if !res {
		fmt.Printf("Failed ibrowser self check - NumSNPS %d != %d\n", ib.NumSNPS, ib.Block.NumSNPS)
		return res
	}

	res = res && (ib.NumSamples == ib.Block.NumSamples)

	if !res {
		fmt.Printf("Failed ibrowser self check - NumSamples %d != %d\n", ib.NumSamples, ib.Block.NumSamples)
		return res
	}

	return res
}

//
// Getters
//

// GetSamples returns the names of the samples
func (ib *IBrowser) GetSamples() VCFSamples {
	return ib.Samples
}

// HasSample checks whether a sample exists
func (ib *IBrowser) HasSample(sampleName string) bool {
	samples := ib.GetSamples()
	_, ok := SliceIndex(len(samples), func(i int) bool { return samples[i] == sampleName })
	return ok
}

// GetSampleID returns the sample ID for a given sample name
func (ib *IBrowser) GetSampleID(sampleName string) (int, bool) {
	samples := ib.GetSamples()
	ind, ok := SliceIndex(len(samples), func(i int) bool { return samples[i] == sampleName })
	return ind, ok
}

// GetSampleName returns the sample name for a given sample ID
func (ib *IBrowser) GetSampleName(sampleID int) (string, bool) {
	samples := ib.GetSamples()
	if sampleID >= len(samples) {
		return "", false
	}
	return samples[sampleID], true
}

// GetSummaryBlock returns the summary block
func (ib *IBrowser) GetSummaryBlock() (*IBBlock, bool) {
	return ib.Block, true
}

// GetSummaryBlockMatrix returns the summary block matrix
func (ib *IBrowser) GetSummaryBlockMatrix() (*IBBlock, *IBDistanceMatrix, bool) {
	block, hasBlock := ib.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, hasMatrix
	}

	return block, matrix, true
}

// GetSummaryBlockMatrixTable returns the summary block matrix table
func (ib *IBrowser) GetSummaryBlockMatrixTable() (*IBBlock, *IBDistanceMatrix, *IBDistanceTable, bool) {
	block, matrix, hasMatrix := ib.GetSummaryBlockMatrix()

	if !hasMatrix {
		return nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, nil, nil, hasTable
	}

	return block, matrix, table, true
}

// GetChromosomeNames returns the name of the chromosomes
func (ib *IBrowser) GetChromosomeNames() (chromosomes []string) {
	numChromosomes := len(ib.ChromosomesNames)
	chromosomes = make([]string, numChromosomes, numChromosomes)

	for cl, chromNamePosPair := range ib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		chromosomes[cl] = chromName
	}

	return
}

// GetChromosomes returns the chromosome instances
func (ib *IBrowser) GetChromosomes() (chromosomes []*IBChromosome) {
	numChromosomes := len(ib.ChromosomesNames)
	chromosomes = make([]*IBChromosome, numChromosomes, numChromosomes)

	for cl, chromNamePosPair := range ib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		chromosomes[cl] = ib.Chromosomes[chromName]
	}

	return
}

// GetChromosome returns a given chromosome by its name
func (ib *IBrowser) GetChromosome(chromosomeName string) (*IBChromosome, bool) {
	if chromosome, ok := ib.Chromosomes[chromosomeName]; ok {
		// fmt.Println("GetChromosome", chromosomeName, "exists", &chromosome)
		return chromosome, ok
	}
	// fmt.Println("GetChromosome", chromosomeName, "DOES NOT exists")
	return nil, false
}

// GetChromosomeSummaryBlock returns the summary block of a given chromosome
func (ib *IBrowser) GetChromosomeSummaryBlock(chromosomeName string) (*IBChromosome, *IBBlock, bool) {
	chrom, hasChrom := ib.GetChromosome(chromosomeName)

	if !hasChrom {
		return nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, hasBlock
	}

	return chrom, block, true
}

// getChromosomeSummaryBlockMatrix returns the summary block matrix of a given chromosome
func (ib *IBrowser) getChromosomeSummaryBlockMatrix(chromosomeName string) (*IBChromosome, *IBBlock, *IBDistanceMatrix, bool) {
	chrom, block, hasChrom := ib.GetChromosomeSummaryBlock(chromosomeName)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	matrix, _ := block.GetMatrix()

	return chrom, block, matrix, true
}

// getChromosomeSummaryBlockMatrixTable returns the summary block matrix table of a given chromosome
func (ib *IBrowser) getChromosomeSummaryBlockMatrixTable(chromosomeName string) (*IBChromosome, *IBBlock, *IBDistanceMatrix, *IBDistanceTable, bool) {
	chrom, block, matrix, hasChrom := ib.getChromosomeSummaryBlockMatrix(chromosomeName)

	if !hasChrom {
		return nil, nil, nil, nil, hasChrom
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, nil, nil, nil, hasTable
	}

	return chrom, block, matrix, table, true
}

// GetChromosomeBlocks returns all blocks of a chromosome
func (ib *IBrowser) GetChromosomeBlocks(chromosomeName string) (*IBChromosome, []*IBBlock, bool) {
	chrom, hasChrom := ib.GetChromosome(chromosomeName)

	if !hasChrom {
		return nil, nil, hasChrom
	}

	blocks, hasBlocks := chrom.GetBlocks()

	if !hasBlocks {
		return nil, nil, hasBlocks
	}

	return chrom, blocks, true
}

// GetChromosomeBlock returns a block from a chromosome given the chromosome name and block number
func (ib *IBrowser) GetChromosomeBlock(chromosomeName string, blockNum uint64) (*IBChromosome, *IBBlock, bool) {
	chrom, hasChrom := ib.GetChromosome(chromosomeName)

	if !hasChrom {
		return nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetBlock(blockNum)

	if !hasBlock {
		return nil, nil, hasBlock
	}

	return chrom, block, true
}

func (ib *IBrowser) getChromosomeBlockMatrix(chromosomeName string, blockNum uint64) (*IBChromosome, *IBBlock, *IBDistanceMatrix, bool) {
	chrom, block, hasBlock := ib.GetChromosomeBlock(chromosomeName, blockNum)

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, hasMatrix
	}

	return chrom, block, matrix, true
}

func (ib *IBrowser) getChromosomeBlockMatrixTable(chromosomeName string, blockNum uint64) (*IBChromosome, *IBBlock, *IBDistanceMatrix, *IBDistanceTable, bool) {
	chrom, block, matrix, hasMatrix := ib.getChromosomeBlockMatrix(chromosomeName, blockNum)

	if !hasMatrix {
		return nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, nil, nil, nil, hasTable
	}

	return chrom, block, matrix, table, true
}

//
// Filename
//

// GenFilename returns the filename of this project when saved
func (ib *IBrowser) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//

// Save saves this project to file
func (ib *IBrowser) Save(outPrefix string, format string, compression string) {
	isSave := true
	isSoft := false
	ib.saveLoad(isSave, isSoft, outPrefix, format, compression)
}

//
// Load
//

// EasyLoadPrefix loads a project from file by its prefix and guesses the file format
func (ib *IBrowser) EasyLoadPrefix(outPrefix string, isSoft bool) {
	found, format, compression, _ := save.GuessPrefixFormat(outPrefix)

	if !found {
		fmt.Println("could not easy load prefix: ", outPrefix)
		os.Exit(1)
	}

	isSave := false
	ib.saveLoad(isSave, isSoft, outPrefix, format, compression)
}

// EasyLoadFile loads a project from file by its full name and guesses the file format
func (ib *IBrowser) EasyLoadFile(outFile string, isSoft bool) {
	found, format, compression, outPrefix := save.GuessFormat(outFile)

	if !found {
		fmt.Println("could not easy load file:", outFile)
		os.Exit(1)
	}

	isSave := false
	ib.saveLoad(isSave, isSoft, outPrefix, format, compression)
}

// Load loads a project from file
func (ib *IBrowser) Load(outPrefix string, format string, compression string, isSoft bool) {
	isSave := false
	ib.saveLoad(isSave, isSoft, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ib *IBrowser) saveLoad(isSave bool, isSoft bool, outPrefix string, format string, compression string) {
	baseName, _ := ib.GenFilename(outPrefix, format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving global ibrowser status")
		ib.Dump(outPrefix, isSave, isSoft)
		saver.Save(ib)
	} else {
		fmt.Println("loading global ibrowser status")
		saver.Load(ib)
		sort.Sort(ib.ChromosomesNames)
		ib.Dump(outPrefix, isSave, isSoft)
	}
}

//
// Dumper
//

// GenMatrixDumpFileName generates the filename of a dump file
func (ib *IBrowser) GenMatrixDumpFileName(outPrefix string) (filename string) {
	filename = outPrefix + "_summary.bin"
	return
}

// GenMatrixChromosomeDumpFileName generates the filename of a dump file for a chromosome
func (ib *IBrowser) GenMatrixChromosomeDumpFileName(outPrefix string, chromosomeName string) (filename string) {
	filename = outPrefix + "_chromosomes_" + chromosomeName + ".bin"
	return
}

// Dump dumps matrices to file
func (ib *IBrowser) Dump(outPrefix string, isSave bool, isSoft bool) {
	summaryFileName := ib.GenMatrixDumpFileName(outPrefix)

	if isSave {
		ib.blockManager.Save(summaryFileName)
	} else {
		ib.blockManager.Load(summaryFileName)
	}

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.Chromosomes[chromosomeName.Name]

		chromosome.DumpBlocks(outPrefix, isSave, isSoft)
	}
}
