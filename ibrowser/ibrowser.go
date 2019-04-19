package ibrowser

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"sync/atomic"
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
	Parameters     Parameters
	//
	ChromosomesNames NamePosPairList
	chromosomes      map[string]*IBChromosome
	block            *IBBlock
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
		Parameters:   parameters,
		//
		lastChrom:    "",
		lastPosition: 0,
		//
		ChromosomesNames: make(NamePosPairList, 0, 100),
		chromosomes:      make(map[string]*IBChromosome, 100),
		//
		// block: NewIBBlock("_whole_genome", blockSize, counterBits, 0, 0, 0),
	}

	return &ib
}

func (ib *IBrowser) SetSamples(samples *VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(VCFSamples, numSamples, numSamples)

	ib.NumSamples = uint64(numSamples)
	ib.block = NewIBBlock("_whole_genome", ib.BlockSize, ib.CounterBits, ib.NumSamples, 0, 0)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

func (ib *IBrowser) GetSamples() VCFSamples {
	return ib.Samples
}

func (ib *IBrowser) GetChromosome(chromosomeName string) (*IBChromosome, bool) {
	if chromosome, ok := ib.chromosomes[chromosomeName]; ok {
		// fmt.Println("GetChromosome", chromosomeName, "exists", &chromosome)
		return chromosome, ok
	} else {
		// fmt.Println("GetChromosome", chromosomeName, "DOES NOT exists")
		return &IBChromosome{}, ok
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

	ib.chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, chromosomeNumber, ib.BlockSize, ib.CounterBits, ib.NumSamples, ib.KeepEmptyBlock)

	ib.ChromosomesNames = append(ib.ChromosomesNames, NamePosPair{chromosomeName, chromosomeNumber})

	sort.Sort(ib.ChromosomesNames)

	return ib.chromosomes[chromosomeName]
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

		ib.block.Add(0, reg.Distance)
	}
	mutex.Unlock()
}

func (ib *IBrowser) Check() (res bool) {
	res = true

	res = res && ib.selfCheck()

	if !res {
		fmt.Println("Failed ibrowser self check")
		return res
	}

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.chromosomes[chromosomeName.Name]

		res = res && chromosome.Check()

		if !res {
			fmt.Printf("Failed ibrowser chromosome %s check\n", chromosomeName)
			return res
		}
	}

	return res
}

func (ib *IBrowser) selfCheck() (res bool) {
	res = true

	res = res && ib.block.Check()

	if !res {
		fmt.Printf("Failed ibrowser self check - block check\n")
		return res
	}

	{
		res = res && (ib.BlockSize == ib.block.BlockSize)

		if !res {
			fmt.Printf("Failed ibrowser self check - block size %d != %d\n", ib.BlockSize, ib.block.BlockSize)
			return res
		}

		res = res && (ib.CounterBits == ib.block.CounterBits)

		if !res {
			fmt.Printf("Failed ibrowser self check - CounterBits %d != %d\n", ib.CounterBits, ib.block.CounterBits)
			return res
		}
	}

	res = res && (ib.NumSNPS == ib.block.NumSNPS)

	if !res {
		fmt.Printf("Failed ibrowser self check - NumSNPS %d != %d\n", ib.NumSNPS, ib.block.NumSNPS)
		return res
	}

	res = res && (ib.NumSamples == ib.block.NumSamples)

	if !res {
		fmt.Printf("Failed ibrowser self check - NumSamples %d != %d\n", ib.NumSamples, ib.block.NumSamples)
		return res
	}

	return res
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
	ib.saveLoad(true, outPrefix, format, compression)
}

//
// Load
//

func (ib *IBrowser) Load(outPrefix string, format string, compression string) {
	ib.saveLoad(false, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ib *IBrowser) dumper(isSave bool, outPrefix string) {
	mode := ""

	if isSave {
		mode = "w"
	} else {
		mode = "r"
	}

	dumper := NewMultiArrayFile(outPrefix+".bin", mode)
	defer dumper.Close()

	ib.dumperMatrix(dumper, isSave, ib.block.GetMatrix())

	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]
		chromosome := ib.chromosomes[chromosomeName.Name]

		ib.dumperMatrix(dumper, isSave, chromosome.block.GetMatrix())

		for _, blockPos := range chromosome.BlockNames {
			block := chromosome.blocks[blockPos]
			ib.dumperMatrix(dumper, isSave, block.GetMatrix())
		}
	}
}

func (ib *IBrowser) dumperMatrix(dumper *MultiArrayFile, isSave bool, data *DistanceMatrix) {
	if isSave {
		if ib.CounterBits == 16 {
			dumper.Write16(&data.Data16)
		} else if ib.CounterBits == 32 {
			dumper.Write32(&data.Data32)
		} else if ib.CounterBits == 64 {
			dumper.Write64(&data.Data64)
		}
	} else {
		if ib.CounterBits == 16 {
			dumper.Read16(&data.Data16)
		} else if ib.CounterBits == 32 {
			dumper.Read32(&data.Data32)
		} else if ib.CounterBits == 64 {
			dumper.Read64(&data.Data64)
		}
	}
}

func (ib *IBrowser) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := ib.GenFilename(outPrefix, format, compression)
	saver := NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving global ibrowser status")
		saver.Save(ib)
	} else {
		fmt.Println("loading global ibrowser status")
		saver.Load(ib)
	}

	ib.saveLoadBlock(isSave, baseName, format, compression)
	ib.saveLoadChromosomes(isSave, baseName, format, compression)

	ib.dumper(isSave, outPrefix)
}

func (ib *IBrowser) saveLoadBlock(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_block"

	if isSave {
		fmt.Println("saving global ibrowser block")
		ib.block.Save(newPrefix, format, compression)
	} else {
		fmt.Println("loading global ibrowser block")
		ib.block = NewIBBlock("_whole_genome", ib.BlockSize, ib.CounterBits, ib.NumSamples, 0, 0)
		ib.block.Load(newPrefix, format, compression)
	}
}

func (ib *IBrowser) saveLoadChromosomes(isSave bool, outPrefix string, format string, compression string) {
	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]

		if isSave {
			fmt.Println("saving chromosome        : ", chromosomeName)
			chromosome := ib.chromosomes[chromosomeName.Name]
			chromosome.Save(outPrefix, format, compression)

		} else {
			fmt.Println("loading chromosome       : ", chromosomeName)
			ib.chromosomes[chromosomeName.Name] = NewIBChromosome(chromosomeName.Name, chromosomeName.Pos, ib.BlockSize, ib.CounterBits, ib.NumSamples, ib.KeepEmptyBlock)
			chromosome := ib.chromosomes[chromosomeName.Name]
			chromosome.Load(outPrefix, format, compression)
		}
	}
}
