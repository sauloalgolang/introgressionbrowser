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

type IBrowser struct {
	Samples        VCFSamples
	NumSamples     uint64
	BlockSize      uint64
	KeepEmptyBlock bool
	NumRegisters   uint64
	NumSNPS        uint64
	NumBlocks      uint64
	CounterBits    int
	RegisterSize   uint64
	Parameters     Parameters
	//
	ChromosomesNames NamePosPairList
	Chromosomes      map[string]*IBChromosome
	Block            *IBBlock
	//
	lastChrom    string
	lastPosition uint64
	//
	// Header string
	//
	// TODO: per sample stats
}

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
	}
	return &ib
}

func (ib *IBrowser) SetSamples(samples *VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(VCFSamples, numSamples, numSamples)

	ib.NumSamples = uint64(numSamples)
	ib.Block = NewIBBlock("_whole_genome", 0, ib.BlockSize, ib.CounterBits, ib.NumSamples, 0, 0)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "exists", &chromosome)
		return chromosome
	} else {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "creating")
		return ib.AddChromosome(chromosomeName, chromosomeNumber)
	}
}

func (ib *IBrowser) AddChromosome(chromosomeName string, chromosomeNumber int) *IBChromosome {
	if chromosome, hasChromosome := ib.GetChromosome(chromosomeName); hasChromosome {
		fmt.Println("Failed to add chromosome", chromosomeName, ". Already exists", &chromosome)
		os.Exit(1)
	}

	ib.Chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, chromosomeNumber, ib.BlockSize, ib.CounterBits, ib.NumSamples, ib.KeepEmptyBlock)

	ib.ChromosomesNames = append(ib.ChromosomesNames, NamePosPair{chromosomeName, chromosomeNumber})

	sort.Sort(ib.ChromosomesNames)

	return ib.Chromosomes[chromosomeName]
}

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

func (ib *IBrowser) GetSamples() VCFSamples {
	return ib.Samples
}

func (ib *IBrowser) HasSample(sampleName string) bool {
	samples := ib.GetSamples()
	_, ok := SliceIndex(len(samples), func(i int) bool { return samples[i] == sampleName })
	return ok
}

func (ib *IBrowser) GetSampleId(sampleName string) (int, bool) {
	samples := ib.GetSamples()
	ind, ok := SliceIndex(len(samples), func(i int) bool { return samples[i] == sampleName })
	return ind, ok
}

func (ib *IBrowser) GetSampleName(sampleId int) (string, bool) {
	samples := ib.GetSamples()
	if sampleId >= len(samples) {
		return "", false
	}
	return samples[sampleId], true
}

func (ib *IBrowser) GetSummaryBlock() (*IBBlock, bool) {
	return ib.Block, true
}

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

func (ib *IBrowser) GetSummaryBlockMatrixData() (*IBBlock, *IBDistanceMatrix, *IBDistanceTable, bool) {
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

func (ib *IBrowser) GetChromosomeNames() (chromosomes []string) {
	numChromosomes := len(ib.ChromosomesNames)
	chromosomes = make([]string, numChromosomes, numChromosomes)

	for cl, chromNamePosPair := range ib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		chromosomes[cl] = chromName
	}

	return
}

func (ib *IBrowser) GetChromosomes() (chromosomes []*IBChromosome) {
	numChromosomes := len(ib.ChromosomesNames)
	chromosomes = make([]*IBChromosome, numChromosomes, numChromosomes)

	for cl, chromNamePosPair := range ib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		chromosomes[cl] = ib.Chromosomes[chromName]
	}

	return
}

func (ib *IBrowser) GetChromosome(chromosomeName string) (*IBChromosome, bool) {
	if chromosome, ok := ib.Chromosomes[chromosomeName]; ok {
		// fmt.Println("GetChromosome", chromosomeName, "exists", &chromosome)
		return chromosome, ok
	} else {
		// fmt.Println("GetChromosome", chromosomeName, "DOES NOT exists")
		return nil, ok
	}
}

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

func (ib *IBrowser) getChromosomeSummaryBlockMatrix(chromosomeName string) (*IBChromosome, *IBBlock, *IBDistanceMatrix, bool) {
	chrom, block, hasChrom := ib.GetChromosomeSummaryBlock(chromosomeName)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	matrix, _ := block.GetMatrix()

	return chrom, block, matrix, true
}

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

func (ib *IBrowser) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix

	saver := NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//
func (ib *IBrowser) Save(outPrefix string, format string, compression string) {
	ib.saveLoad(true, outPrefix, format, compression, false)
}

//
// Load
//

func (ib *IBrowser) EasyLoadPrefix(outPrefix string, soft bool) {
	found, format, compression, _ := save.GuessPrefixFormat(outPrefix)

	if !found {
		fmt.Println("could not easy load prefix: ", outPrefix)
		os.Exit(1)
	}

	ib.saveLoad(false, outPrefix, format, compression, soft)
}

func (ib *IBrowser) EasyLoadFile(outFile string, soft bool) {
	found, format, compression, outPrefix := save.GuessFormat(outFile)

	if !found {
		fmt.Println("could not easy load file:", outFile)
		os.Exit(1)
	}

	ib.saveLoad(false, outPrefix, format, compression, soft)
}

func (ib *IBrowser) Load(outPrefix string, format string, compression string, soft bool) {
	ib.saveLoad(false, outPrefix, format, compression, soft)
}

//
// SaveLoad
//

func (ib *IBrowser) saveLoad(isSave bool, outPrefix string, format string, compression string, soft bool) {
	baseName, _ := ib.GenFilename(outPrefix, format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving global ibrowser status")
		ib.dumper(isSave, outPrefix)
		saver.Save(ib)
	} else {
		fmt.Println("loading global ibrowser status")
		saver.Load(ib)
		sort.Sort(ib.ChromosomesNames)
		if !soft {
			ib.dumper(isSave, outPrefix)
		}
	}

	// ib.saveLoadBlock(isSave, baseName, format, compression)
	// ib.saveLoadChromosomes(isSave, baseName, format, compression)
}

func (ib *IBrowser) saveLoadBlock(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_block"

	if isSave {
		fmt.Println("saving global ibrowser block")
		ib.Block.Save(newPrefix, format, compression)
	} else {
		fmt.Println("loading global ibrowser block")
		ib.Block = NewIBBlock(
			"_whole_genome",
			0,
			ib.BlockSize,
			ib.CounterBits,
			ib.NumSamples,
			0,
			0,
		)
		ib.Block.Load(newPrefix, format, compression)
	}
}

func (ib *IBrowser) saveLoadChromosomes(isSave bool, outPrefix string, format string, compression string) {
	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]

		if isSave {
			fmt.Println("saving chromosome        : ", chromosomeName)
			chromosome := ib.Chromosomes[chromosomeName.Name]
			chromosome.Save(outPrefix, format, compression)

		} else {
			fmt.Println("loading chromosome       : ", chromosomeName)
			ib.Chromosomes[chromosomeName.Name] = NewIBChromosome(chromosomeName.Name, chromosomeName.Pos, ib.BlockSize, ib.CounterBits, ib.NumSamples, ib.KeepEmptyBlock)
			chromosome := ib.Chromosomes[chromosomeName.Name]
			chromosome.Load(outPrefix, format, compression)
		}
	}
}

//
// Dumper
//
func (ib *IBrowser) GenMatrixDumpFileName(outPrefix string, chromosomeName string, isSummary bool, isChromosomes bool) (filename string) {
	if isSummary {
		if isChromosomes {
			filename = outPrefix + "_chromosomes.bin"
		} else {
			filename = outPrefix + "_summary.bin"
		}
	} else {
		filename = outPrefix + "_chromosomes_" + chromosomeName + ".bin"
	}
	return
}

func (ib *IBrowser) dumper(isSave bool, outPrefix string) {
	mode := ""

	if isSave {
		mode = "w"
	} else {
		mode = "r"
	}

	summaryFileName := ib.GenMatrixDumpFileName(outPrefix, "", true, false)
	// summaryChromFileName := ib.GenMatrixDumpFileName(outPrefix, "", true, true)

	dumperg := NewMultiArrayFile(summaryFileName, mode)
	// dumperc := NewMultiArrayFile(summaryChromFileName, mode)

	ib.RegisterSize = dumperg.CalculateRegisterSize(ib.CounterBits, ib.Block.Matrix.Size)

	defer dumperg.Close()
	// defer dumperc.Close()

	ib.Block.Dump(dumperg, isSave)
	// ib.dumperMatrix(dumperg, isSave, ib.Block)

	// fmt.Println("ib.ChromosomesNames", ib.ChromosomesNames)

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.Chromosomes[chromosomeName.Name]

		// ib.dumperMatrix(dumperg, isSave, chromosome.Block)
		chromosome.Block.Dump(dumperg, isSave)

		// outPrefix+"_chromosomes_"+chromosomeName.Name+".bin"
		chromosomeFileName := ib.GenMatrixDumpFileName(outPrefix, chromosomeName.Name, false, false)
		dumperl := NewMultiArrayFile(chromosomeFileName, mode)
		// dumperl.SetSerial(dumperc.GetSerial())

		for _, block := range chromosome.Blocks {
			// ib.dumperMatrix(dumperc, isSave, block)
			// ib.dumperMatrix(dumperl, isSave, block)
			block.Dump(dumperl, isSave)
		}

		dumperl.Close()
	}
}
