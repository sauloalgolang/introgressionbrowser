package ibrowser

import (
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/save"
)

const defaultGlobalSummaryName = "_whole_genome"

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
	BlockManager     *BlockManager
	//
	lastChrom    string
	lastPosition uint64
	outPrefix    string
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
	outPrefix := parameters.Outfile

	if blockSize > uint64((math.MaxUint32/3)-1) {
		log.Println("block size too large")
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
		outPrefix: outPrefix,
	}

	return &ib
}

// SetSamples set the sample names
func (ib *IBrowser) SetSamples(samples *VCFSamples) {
	numSamples := len(*samples)

	ib.Samples = make(VCFSamples, numSamples, numSamples)
	ib.NumSamples = uint64(numSamples)

	fileName := ib.GetMatrixDumpFileName()

	ib.BlockManager = NewBlockManager(
		defaultGlobalSummaryName,
		fileName,
		ib.CounterBits,
		ib.NumSamples,
		ib.BlockSize,
	)

	chromosomeNumber := int(0)
	blockNumber := uint64(0)
	ib.BlockManager.NewBlock(
		defaultGlobalSummaryName,
		chromosomeNumber,
		blockNumber,
	)

	for samplePos, sampleName := range *samples {
		// log.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

// GetOrCreateChromosome gets or creates a new chromosome if it does not already exists
func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		// log.Println("GetOrCreateChromosome", chromosomeName, "exists", &chromosome)
		return chromosome
	}
	// log.Println("GetOrCreateChromosome", chromosomeName, "creating")
	return ib.AddChromosome(chromosomeName, chromosomeNumber)
}

// AddChromosome adds a new chromosome
func (ib *IBrowser) AddChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, hasChromosome := ib.GetChromosome(chromosomeName); hasChromosome {
		log.Println("Failed to add chromosome", chromosomeName, ". Already exists", &chromosome)
		os.Exit(1)
	}

	fileName := ib.GetMatrixChromosomeDumpFileName(chromosomeName)

	ib.Chromosomes[chromosomeName] = NewIBChromosome(
		chromosomeName,
		chromosomeNumber,
		ib.BlockSize,
		ib.CounterBits,
		ib.NumSamples,
		ib.KeepEmptyBlock,
		fileName,
		ib.BlockManager,
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
			log.Println("Sample mismatch")
			log.Println(len(ib.Samples), "!=", len(*samples))
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

		block, hasBlock := ib.GetSummaryBlock()
		if !hasBlock {
			log.Println(ib.BlockManager)
			panic("!GetSummaryBlock")
		}

		block.AddVcfMatrix(0, reg.Distance)
	}
	mutex.Unlock()
}

// Check checks for self consistency in the data
func (ib *IBrowser) Check() (res bool) {
	log.Println("Starting self check")

	res = true

	res = res && ib.selfCheck()

	if !res {
		log.Println("Failed ibrowser self check")
		return res
	}

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.Chromosomes[chromosomeName.Name]

		res = res && chromosome.Check()

		if !res {
			log.Printf("Failed ibrowser chromosome %s (%d) check\n", chromosomeName.Name, chromosomeName.Pos)
			return res
		}
	}

	return res
}

func (ib *IBrowser) selfCheck() (res bool) {
	res = true

	block, hasBlock := ib.GetSummaryBlock()
	if !hasBlock {
		log.Println(ib.BlockManager)
		panic("!GetSummaryBlock")
	}

	res = res && block.Check()

	if !res {
		log.Printf("Failed ibrowser self check - block check\n")
		return res
	}

	{
		res = res && (ib.BlockSize == block.BlockSize)

		if !res {
			log.Printf("Failed ibrowser self check - block size %d != %d\n", ib.BlockSize, block.BlockSize)
			return res
		}

		res = res && (ib.CounterBits == block.CounterBits)

		if !res {
			log.Printf("Failed ibrowser self check - CounterBits %d != %d\n", ib.CounterBits, block.CounterBits)
			return res
		}
	}

	res = res && (ib.NumSNPS == block.NumSNPS)

	if !res {
		log.Printf("Failed ibrowser self check - NumSNPS %d != %d\n", ib.NumSNPS, block.NumSNPS)
		return res
	}

	res = res && (ib.NumSamples == block.NumSamples)

	if !res {
		log.Printf("Failed ibrowser self check - NumSamples %d != %d\n", ib.NumSamples, block.NumSamples)
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
func (ib *IBrowser) GetSummaryBlock() (block *IBBlock, hasBlock bool) {
	block, hasBlock = ib.BlockManager.GetBlockByName(defaultGlobalSummaryName)
	if !hasBlock {
		log.Println(ib.BlockManager)
		panic("!GetSummaryBlock")
	}
	return block, hasBlock
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
		// log.Println("GetChromosome", chromosomeName, "exists", &chromosome)
		return chromosome, ok
	}
	// log.Println("GetChromosome", chromosomeName, "DOES NOT exists")
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
func (ib *IBrowser) GenFilename(format string, compression string) (baseName string, fileName string) {
	baseName = ib.outPrefix

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//

// Save saves this project to file
func (ib *IBrowser) Save(format string, compression string) {
	isSave := true
	ib.saveLoad(isSave, format, compression)
}

//
// Load
//

// EasyLoadPrefix loads a project from file by its prefix and guesses the file format
func (ib *IBrowser) EasyLoadPrefix() {
	found, format, compression, _ := save.GuessPrefixFormat(ib.outPrefix)

	if !found {
		log.Println("could not easy load prefix: ", ib.outPrefix)
		os.Exit(1)
	}

	isSave := false
	ib.saveLoad(isSave, format, compression)
}

// EasyLoadFile loads a project from file by its full name and guesses the file format
func (ib *IBrowser) EasyLoadFile(outFile string) {
	found, format, compression, outPrefix := save.GuessFormat(outFile)

	if !found {
		log.Println("could not easy load file:", outFile)
		os.Exit(1)
	}

	ib.outPrefix = outPrefix

	isSave := false
	ib.saveLoad(isSave, format, compression)
}

// Load loads a project from file
func (ib *IBrowser) Load(format string, compression string) {
	isSave := false
	ib.saveLoad(isSave, format, compression)
}

//
// SaveLoad
//

func (ib *IBrowser) saveLoad(isSave bool, format string, compression string) {
	baseName, _ := ib.GenFilename(format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		ib.Dump(isSave)
		log.Println("saving global ibrowser status")
		saver.Save(ib)
		log.Println("saving global ibrowser status - DONE")
	} else {
		log.Println("loading global ibrowser status")
		saver.Load(ib)
		sort.Sort(ib.ChromosomesNames)
		log.Println("loading global ibrowser status - DONE")
		ib.Dump(isSave)
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

// GetMatrixDumpFileName generates the filename of a dump file
func (ib *IBrowser) GetMatrixDumpFileName() (string) {
	return ib.GenMatrixDumpFileName(ib.outPrefix)
}

// GetMatrixChromosomeDumpFileName generates the filename of a dump file for a chromosome
func (ib *IBrowser) GetMatrixChromosomeDumpFileName(chromosomeName string) (string) {
	return ib.GenMatrixChromosomeDumpFileName(ib.outPrefix, chromosomeName)
}

// Dump dumps matrices to file
func (ib *IBrowser) Dump(isSave bool) {
	// summaryFileName := ib.GenMatrixDumpFileName()

	// log.Println("(Un)Dumping Matrices")

	// if isSave {
	// 	log.Println(" Dumping Summary Matrices")
	// 	// ib.BlockManager.Save(summaryFileName)
	// 	log.Println(" Dumping Summary Matrices - DONE")
	// } else {
	// 	log.Println(" UnDumping Summary Matrices")
	// 	// ib.BlockManager.Load(summaryFileName)
	// 	log.Println(" UnDumping Summary Matrices - DONE")
	// }

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.Chromosomes[chromosomeName.Name]

		chromosome.setRootBlockManager(ib.BlockManager)
		// chromosome.DumpBlocks(outPrefix, isSave, isSoft)
	}

	log.Println("(Un) Dumping Matrices - DONE")
}
